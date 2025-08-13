package handlers

import (
	"go-EdTech/config"
	"go-EdTech/logger"
	"go-EdTech/models"
	"go-EdTech/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	usersRepo *repositories.UsersRepository
}

func NewAuthHandlers(usersRepo *repositories.UsersRepository) *AuthHandlers {
	return &AuthHandlers{usersRepo: usersRepo}
}

type signinRequest struct {
	Email    string
	Password string
}

type registerNewUserRequest struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	User_Type string `json:"user_type"`
}

// SignIn godoc
// @Summary 	User SignIn func using email and password
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param		email 		query 		string 		true 		"Email"
// @Param 		password 	query 		string 		true 		"Password"
// @Success 	200 		{object} 	object{token=string} 	"OK"
// @Failure 	400 		{object} 	models.ApiError 		"Invalid request Payload"
// @Failure 	401 		{object} 	models.ApiError 		"Invalid credentials"
// @Failure 	500 		{object} 	models.ApiError
// @Router 		/auth/signIn [post]
func (h *AuthHandlers) SignIn(c *gin.Context) {
	var request signinRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}
	user, err := h.usersRepo.FindByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewApiError("Invalid credentials"))
		return
	}

	claims := jwt.RegisteredClaims{
		Subject:   user.UUID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.JwtExpiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not generate token"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

// SignOut godoc
// @Summary 	Sign out
// @Tags 		users
// @Accept 		json
// @Produce		json
// @Success 	200 		"OK"
// @Router 		/auth/singOut [post]
// @Security Bearer
func (h *AuthHandlers) SignOut(c *gin.Context) {
	c.Status(http.StatusOK)
}

// GetUserInfo godoc
// @Summary 	Get user Info
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Success 	200 	{object} userResponse 		"OK"
// @Failure 	500 	{object} models.ApiError
// @Router 		/auth/userInfo [get]
// @Security Bearer
func (h *AuthHandlers) GetUserInfo(c *gin.Context) {
	userUUId := c.GetString("userUUId")
	user, err := h.usersRepo.FindById(c, userUUId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, userResponse{
		UUID:      user.UUID,
		Email:     user.Email,
		Name:      user.Name,
		Surname:   user.Surname,
		User_Type: user.User_Type,
	})
}

// Create godoc
// @Summary 	Create User
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param		name 		query 		string 		true 		"Name"
// @Param		surname 	query 		string 		true 		"Surname"
// @Param		email 		query 		string 		true 		"Email"
// @Param		type 		query 		string 		true 		"User_Type" Enum('Child', 'Parent', 'Content_manager', 'Administrator')
// @Param		password 	query 		string 		true 		"Password"
// @Success 	200 		{object} 	object{uuid=string}		"OK"
// @Failure 	400 		{object} 	models.ApiError 		"Invalid Payload"
// @Failure 	500 		{object} 	models.ApiError
// @Router 		/auth/signUp [post]
// @Security Bearer
func (h *UsersHandlers) RegisterNewUser(c *gin.Context) {
	logger := logger.GetLogger()
	var request registerNewUserRequest
	err := c.BindJSON(&request)
	if err != nil {
		logger.Error("Failed JSON binding", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid payload"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error with password encrypting", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("fail to hash password"))
		return
	}

	user := models.User{
		Name:         request.Name,
		Surname:      request.Surname,
		Email:        request.Email,
		User_Type:    request.User_Type,
		PasswordHash: string(passwordHash),
	}

	id, err := h.repo.Create(c, user)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Could not create user"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"uuid": id})

}

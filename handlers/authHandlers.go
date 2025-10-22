package handlers

import (
	"context"
	"go-EdTech/config"
	"go-EdTech/logger"
	"go-EdTech/models"
	"go-EdTech/repositories"
	"go-EdTech/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	usersRepo    *repositories.UsersRepository
	sessionsRepo *repositories.SessionsRepository
	rolesRepo    *repositories.RoleRepository
}

func NewAuthHandler(usersRepo *repositories.UsersRepository, sessionsRepo *repositories.SessionsRepository, rolesRepo *repositories.RoleRepository) *AuthHandlers {
	return &AuthHandlers{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
		rolesRepo:    rolesRepo,
	}
}

type signinRequest struct {
	Email    string
	Password string
}

type registerNewUserRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role_id  int    `json:"role_id"`
}

// Login godoc
// @Summary 	Вход пользователя в аккаунт, используя email и пароль
// @Tags 		Auth
// @Accept 		json
// @Produce 	json
// @Param		email 		query 		string 		true 		"Email"
// @Param 		password 	query 		string 		true 		"Password"
// @Success 	200 		{object} 	object{token=string, user=models.User, role=string} 	"OK"
// @Failure 	400 		{object} 	models.ApiError 		"Invalid request Payload"
// @Failure 	401 		{object} 	models.ApiError 		"Invalid credentials"
// @Failure 	500 		{object} 	models.ApiError
// @Router 		/auth/login [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	logger := logger.GetLogger()

	var request signinRequest
	if err := c.BindJSON(&request); err != nil {
		logger.Error("Failed to bind json", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}
	user, err := h.usersRepo.FindByEmail(c, request.Email)
	if err != nil {
		logger.Error("Failed to find user", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	role, err := h.rolesRepo.GetRoleByID(c, user.Role_id)
	if err != nil {
		logger.Error("Failed to find user role", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		logger.Error("Failed to succesfully compare passwords", zap.Error(err))
		c.JSON(http.StatusUnauthorized, models.NewApiError("Invalid credentials"))
		return
	}

	token, err := h.generateJWTToken(c.Request.Context(), user.UUID.String(), user.Role_id)
	if err != nil {
		logger.Error("Failed to generate JWT", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("failed to generate token"))
		return
	}

	// Генерация refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.UUID.String())
	if err != nil {
		logger.Error("Failed to generate refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("failed to generate refresh token"))
		return
	}

	// Создаем сессию
	session := models.Session{
		UserUUID:     user.UUID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(config.Config.TokenExpirationDate),
	}
	if err := h.sessionsRepo.CreateSession(c.Request.Context(), session); err != nil {
		logger.Error("Failed to create session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("failed to create session"))
		return
	}

	// Ставим cookie
	c.SetCookie("session_token", refreshToken, int(session.ExpiresAt.Unix()), "/", "", false, true)

	logger.Info("Application login success", zap.String("user_uuid", user.UUID.String()), zap.String("role", role.Name))

	// Отдаем ответ
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
		"role":  role.Name,
	})

}

// Logout godoc
// @Summary Выход из системы
// @Description Завершает текущую сессию пользователя
// @Tags Auth
// @Produce json
// @Success 200 {object} object{message=string}
// @Failure 400 {object} models.ApiError
// @Failure 500 {object} models.ApiError
// @Security ApiKeyAuth
// @Router /auth/logout [post]
func (h *AuthHandlers) Logout(c *gin.Context) {
	logger := logger.GetLogger()
	// Получаем session token из cookie
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		logger.Warn("Logout attempt without session token")
		c.JSON(http.StatusBadRequest, models.NewApiError("no session token"))
		return
	}

	// Удаляем сессию по session token
	if err := h.sessionsRepo.DeleteSession(c.Request.Context(), sessionToken); err != nil {
		logger.Error("Failed to delete session", zap.String("session_token", sessionToken), zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("failed to delete session"))
		return
	}

	// Удаляем cookie с session token
	c.SetCookie("session_token", "", -1, "/", "", false, true)

	logger.Info("Successful logout", zap.String("session_token", sessionToken))

	// Ответ о успешном выходе
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

func (h *AuthHandlers) generateJWTToken(c context.Context, userUUID string, roleID int) (string, error) {
	logger := logger.GetLogger()
	// Находим пользователя по его ID
	user, err := h.usersRepo.FindByUUID(c, userUUID)
	if err != nil {
		logger.Error("Failed to find user by ID",
			zap.String("user_UUID", userUUID),
			zap.Error(err))
		return "", err
	}

	// Получаем роль пользователя
	role, err := h.rolesRepo.GetRoleByID(c, user.Role_id)
	if err != nil {
		logger.Error("Failed to find user by ID",
			zap.String("user_UUID", userUUID),
			zap.Error(err))
		return "", err
	}

	// Создаем JWT токен с ролью и ID пользователя
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     userUUID,
		"role":    role.Name,
		"role_id": roleID,                                // Добавлено role_id для более удобной проверки
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Время истечения токена — 24 часа
	})

	logger.Debug("JWT token generated", zap.String("userID", userUUID), zap.String("role", role.Name))

	// Подписываем и возвращаем токен
	return token.SignedString([]byte(config.Config.JwtSecretKey))
}

// Sign up godoc
// @Summary 	Создание пользователя. Регистрация
// @Tags 		Auth
// @Accept 		json
// @Produce 	json
// @Param		name 		query 		string 		true 													"Name"
// @Param		surname 	query 		string 		true 													"Surname"
// @Param		email 		query 		string 		true 													"Email"
// @Param		role_id 	query 		int 		true 													"Role_id roles are: 1-Child, 2-Parent, 3-Content-manager, 4-Administrator"
// @Param		password 	query 		string 		true 													"Password"
// @Success 	200 		{object} 	object{uuid=string, token=string, user=models.User, role=string}	"OK"
// @Failure 	400 		{object} 	models.ApiError 		"Invalid Payload"
// @Failure 	500 		{object} 	models.ApiError
// @Router 		/auth/signup [post]
// @Security ApiKeyAuth
func (h *AuthHandlers) SignUp(c *gin.Context) {
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
		Role_id:      request.Role_id,
		PasswordHash: string(passwordHash),
	}

	id, err := h.usersRepo.Create(c, user)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Could not create user"))
		return
	}

	user.UUID = id

	// Получаем роль пользователя
	role, err := h.rolesRepo.GetRoleByID(c, user.Role_id)
	if err != nil {
		logger.Error("Failed to get role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Could not find role"))
		return
	}

	// Генерация JWT
	token, err := h.generateJWTToken(c.Request.Context(), id.String(), user.Role_id)
	if err != nil {
		logger.Error("Failed to generate JWT", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("failed to generate token"))
		return
	}

	logger.Info("Creation success", zap.String("user_uuid", user.UUID.String()), zap.String("role", role.Name))

	// Отдаем ответ
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
		"role":  role.Name,
	})

}

// Refresh godoc
// @Summary Обновление токена
// @Description Обновляет JWT токен с помощью refresh токена из cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} object{token=string,expires=int}
// @Failure 401 {object} models.ApiError
// @Failure 500 {object} models.ApiError
// @Router /auth/refresh [post]
func (h *AuthHandlers) Refresh(c *gin.Context) {
	logger := logger.GetLogger()
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		logger.Warn("Refresh attempt without session token")
		c.JSON(http.StatusUnauthorized, models.NewApiError("no session token"))
		return
	}

	session, roleID, err := h.sessionsRepo.GetSession(c.Request.Context(), sessionToken)
	if err != nil {
		logger.Warn("Invalid session token", zap.String("session_token", sessionToken), zap.Error(err))
		c.JSON(http.StatusUnauthorized, models.NewApiError("invalid session token"))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		logger.Warn("Expired session token", zap.String("session_token", sessionToken))
		c.JSON(http.StatusUnauthorized, models.NewApiError("expired session token"))
		return
	}

	token, err := h.generateJWTToken(c.Request.Context(), session.UserUUID.String(), roleID)
	if err != nil {
		logger.Error("Failed to generate JWT token",
			zap.String("user_uuid", session.UserUUID.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("failed to generate token"))
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(session.UserUUID.String())
	if err != nil {
		logger.Error("Failed to generate refresh token",
			zap.String("user_uuid", session.UserUUID.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("failed to generate refresh token"))
		return
	}

	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(config.Config.TokenExpirationDate)

	if err := h.sessionsRepo.UpdateSession(c.Request.Context(), session); err != nil {
		logger.Error("Failed to update session",
			zap.String("user_uuid", session.UserUUID.String()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("failed to update session"))
		return
	}

	c.SetCookie("session_token", newRefreshToken, int(session.ExpiresAt.Unix()), "/", "", false, true)

	logger.Info("Tokens refreshed successfully", zap.String("user_uuid", session.UserUUID.String()))

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"expires": time.Now().Add(time.Hour * 24).Unix(),
	})
}

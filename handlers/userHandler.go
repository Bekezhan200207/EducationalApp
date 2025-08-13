package handlers

import (
	"go-EdTech/logger"
	"go-EdTech/models"
	"go-EdTech/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandlers struct {
	repo *repositories.UsersRepository
}

func NewUsersHandlers(repo *repositories.UsersRepository) *UsersHandlers {
	return &UsersHandlers{repo: repo}
}

type userResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	User_Type string    `json:"user_type"`
}

type changePasswordRequest struct {
	Password string `json:"password"`
}

// FindAll godoc
// @Summary 	Get All Users
// @Tags 		users
// @Accept	 	json
// @Produce 	json
// @Success 	200 	{object} []userResponse "OK"
// @Failure 	500 	{object} models.ApiError
// @Router 		/users [get]
// @Security Bearer
func (h *UsersHandlers) FindAll(c *gin.Context) {
	logger := logger.GetLogger()

	users, err := h.repo.FindAll(c)
	if err != nil {
		logger.Error("Failed to fetch all users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not load users"))
		return
	}

	dtos := make([]userResponse, 0, len(users))
	for _, u := range users {
		r := userResponse{
			UUID:      u.UUID,
			Name:      u.Name,
			Surname:   u.Surname,
			Email:     u.Email,
			User_Type: u.User_Type,
		}
		dtos = append(dtos, r)
	}

	c.JSON(http.StatusOK, dtos)
}

// FindById godoc
// @Summary 	Find By Id
// @Tags		users
// @Accept 		json
// @Produce 	json
// @Param 		uuid 	path 		string 	true 	"User UUID"
// @Success 	200 	{object} 	userResponse 	"OK"
// @Failure 	500 	{object} 	models.ApiError "Invalid user uuid"
// @Router 		/users/{uuid} [get]
// @Security Bearer
func (h *UsersHandlers) FindById(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")
	user, err := h.repo.FindById(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	dto := userResponse{
		UUID:      user.UUID,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		User_Type: user.User_Type,
	}

	c.JSON(http.StatusOK, dto)
}

// FindByEmail godoc
// @Summary 	Find By Email
// @Tags		users
// @Accept 		json
// @Produce 	json
// @Param 		email 	path 		string 	true 	"Email"
// @Success 	200 	{object} 	userResponse 	"OK"
// @Failure 	400 	{object} 	models.ApiError "Invalid user uuid"
// @Router 		/users/{email} [get]
// @Security Bearer
func (h *UsersHandlers) FindByEmail(c *gin.Context) {
	logger := logger.GetLogger()
	email := c.Param("email")
	user, err := h.repo.FindByEmail(c, email)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	dto := userResponse{
		UUID:      user.UUID,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		User_Type: user.User_Type,
	}

	c.JSON(http.StatusOK, dto)
}

// Update godoc
// @Summary 	Update User
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		uuid 		path 		string 		true 	"User UUID"
// @Param 		name 		query		string 		true 	"User Name"
// @Param 		surname 	query		string 		true 	"User Surname"
// @Param 		email 		query		string 		true 	"Email"
// @Success 	200 		"OK"
// @Failure 	500 		{object}	models.ApiError 	"Invalid user uuid"
// @Router 		/users/{uuid} [put]
// @Security Bearer
func (h *UsersHandlers) Update(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindById(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Invalid user id"))
		return
	}

	var userUpdateRequest models.User
	err = c.BindJSON(&userUpdateRequest)
	if err != nil {
		logger.Error("Failed JSON binding", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not bind Json"))
		return
	}

	err = h.repo.Update(c, userUpdateRequest, uuid)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// ChangePassword godoc
// @Summary 	ChangePassword User
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		uuid 		path 		string 		true 	"User UUID"
// @Param 		password 	query 		string	 	true 	"Password"
// @Success 	200 		"OK"
// @Failure 	500 		{object} 	models.ApiError	 	"Invalid user uuid"
// @Router 		/users/{uuid}/changePassword [patch]
// @Security Bearer
func (h *UsersHandlers) ChangePassword(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindById(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Invalid user id"))
		return
	}

	var userPasswordUpdate changePasswordRequest
	err = c.BindJSON(&userPasswordUpdate)
	if err != nil {
		logger.Error("Failed JSON binding", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not bind Json"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userPasswordUpdate.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed password encryption", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("problem with Hashing"))
		return
	}

	err = h.repo.ChangePassword(c, passwordHash, uuid)
	if err != nil {
		logger.Error("Failed to change user password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary 	Delete User
// @Tags 		users
// @Accept		json
// @Produce 	json
// @Param 		uuid 	path 		string 		true 	"User UUID"
// @Success 	200 	"OK"
// @Failure 	500 	{object} 	models.ApiError 	"Invalid user uuid"
// @Router 		/users/{uuid} [delete]
// @Security Bearer
func (h *UsersHandlers) Delete(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindById(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Invalid user id"))
		return
	}

	err = h.repo.Delete(c, uuid)
	if err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Deactivate godoc
// @Summary 	set user status to inactive
// @Tags 		users
// @Accept 		json
// @Produce		json
// @Param 		uuid 		path 		string 		true 	"User UUID"
// @Success 	200 		"OK"
// @Failure 	500 		{object} 	models.ApiError 	"Invalid user uuid"
// @Router 		/users/{uuid}/deactivate [patch]
// @Security Bearer
func (h *UsersHandlers) Deactivate(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindById(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Invalid user uuid"))
		return
	}

	err = h.repo.Deactivate(c, uuid)
	if err != nil {
		logger.Error("Failed to deactivate user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// Activate godoc
// @Summary 	set user status to active
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		uuid 	path 		string 		true 	"User UUID"
// @Success 	200 	"OK"
// @Failure 	500 	{object} 	models.ApiError 	"Invalid user uuid"
// @Router 		/users/{uuid}/activate [patch]
// @Security Bearer
func (h *UsersHandlers) Activate(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindById(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Invalid user id"))
		return
	}

	err = h.repo.Activate(c, uuid)
	if err != nil {
		logger.Error("Failed to activate user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

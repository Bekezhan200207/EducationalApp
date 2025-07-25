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

type createUserRequest struct {
	User_Name    string
	User_Surname string
	Email        string
	Password     string
	User_Type    string
}

type userResponse struct {
	Id           uuid.UUID `json:"uuid"`
	User_Name    string    `json:"name"`
	User_Surname string    `json:"surname"`
	Email        string    `json:"email"`
	User_Type    string    `json:"user_type"`
}

// FindAll godoc
// @Summary Get All Users
// @Tags users
// @Accept json
// @Produce json
// @Success 		200 {object} []userResponse "OK"
// @Failure 		500 {object} models.ApiError
// @Router 		/users [get]
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
			Id:           u.Id,
			User_Name:    u.User_Name,
			User_Surname: u.User_Surname,
			Email:        u.Email,
			User_Type:    u.User_Type,
		}
		dtos = append(dtos, r)
	}

	c.JSON(http.StatusOK, dtos)
}

// FindOne godoc
// @Summary Find By Id
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Success 		200 {object} userResponse "OK"
// @Failure 		400 {object} models.ApiError "Invalid user uuid"
// @Failure 		500 {object} models.ApiError
// @Router 		/users/{uuid} [get]
// @Security Bearer
func (h *UsersHandlers) FindOne(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")
	user, err := h.repo.FindOne(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	dto := userResponse{
		Id:           user.Id,
		User_Name:    user.User_Name,
		User_Surname: user.User_Surname,
		Email:        user.Email,
		User_Type:    user.User_Type,
	}

	c.JSON(http.StatusOK, dto)
}

// Create godoc
// @Summary Create User
// @Tags users
// @Accept json
// @Produce json
// @Param name query string true "User_Name"
// @Param surname query string true "User_Surname"
// @Param email query string true "Email"
// @Param type query string true "User_Type" Enum('Child', 'Parent', 'Content_manager', 'Administrator')
// @Param password query string true "Password"
// @Success 		200 {object} object{uuid=string} "OK"
// @Failure 		400 {object} models.ApiError "Invalid Payload"
// @Failure 		500 {object} models.ApiError
// @Router 		/users [post]
func (h *UsersHandlers) Create(c *gin.Context) {
	logger := logger.GetLogger()
	var request createUserRequest
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
		User_Name:    request.User_Name,
		User_Surname: request.User_Surname,
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

// Update godoc
// @Summary Update User
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Param name query string true "User Name"
// @Param surname query string true "User Surname"
// @Param email query string true "Email"
// @Success 		200 "OK"
// @Failure 		400 {object} models.ApiError "Invalid user uuid"
// @Failure 		500 {object} models.ApiError
// @Router 		/users/{uuid} [put]
func (h *UsersHandlers) Update(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindOne(c, uuid)
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
// @Summary ChangePassword User
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Param password query string true "Password"
// @Success 		200 "OK"
// @Failure 		400 {object} models.ApiError "Invalid user uuid"
// @Failure 		500 {object} models.ApiError
// @Router 		/users/{uuid}/changePassword [patch]
func (h *UsersHandlers) ChangePassword(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindOne(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Invalid user id"))
		return
	}

	var userPasswordUpdate createUserRequest
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
// @Summary Delete User
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Success 		200 "OK"
// @Failure 		400 {object} models.ApiError "Invalid user uuid"
// @Failure 		500 {object} models.ApiError
// @Router 		/users/{uuid} [delete]
func (h *UsersHandlers) Delete(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindOne(c, uuid)
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
// @Summary set user status to inactive
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Success 		200 "OK"
// @Failure 		400 {object} models.ApiError "Invalid user uuid"
// @Failure 		500 {object} models.ApiError
// @Router 		/users/{uuid}/deactivate [patch]
func (h *UsersHandlers) Deactivate(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindOne(c, uuid)
	if err != nil {
		logger.Error("Requested user not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError("Invalid user id"))
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
// @Summary set user status to active
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Success 		200 "OK"
// @Failure 		400 {object} models.ApiError "Invalid user uuid"
// @Failure 		500 {object} models.ApiError
// @Router 		/users/{uuid}/activate [patch]
func (h *UsersHandlers) Activate(c *gin.Context) {
	logger := logger.GetLogger()
	uuid := c.Param("uuid")

	_, err := h.repo.FindOne(c, uuid)
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

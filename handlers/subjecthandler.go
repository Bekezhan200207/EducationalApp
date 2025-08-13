package handlers

import (
	"go-EdTech/logger"
	"go-EdTech/models"
	"go-EdTech/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubjectsHandlers struct {
	repo *repositories.SubjectsRepository
}

func NewSubjectsHandlers(repo *repositories.SubjectsRepository) *SubjectsHandlers {
	return &SubjectsHandlers{repo: repo}
}

// FindById godoc
// @Summary 	find subject by id
// @Tags 		subjects
// @Accept		json
// @Produce 	json
// @Param 		id 		path 		int 	true 	"Subject_id"
// @Success 	200 	{object} 	models.Subject 	"OK"
// @Failure 	400 	{object} 	models.ApiError "Invalid Payload"
// @Failure 	500 	{object} 	models.ApiError
// @Router 		/subjects/{id} [get]
// @Security 	Bearer
func (g *SubjectsHandlers) FindById(c *gin.Context) {
	logger := logger.GetLogger()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid subject Id format", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Subject not found"))
		return
	}

	subject, err := g.repo.FindById(c, id)
	if err != nil {
		logger.Error("Requested subject not found", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, subject)

}

// FindAll godoc
// @Summary 	find all subjects
// @Tags 		subjects
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} []models.Subject "OK"
// @Failure 	500 {object} models.ApiError
// @Router 		/subjects [get]
// @Security 	Bearer
func (g *SubjectsHandlers) FindAll(c *gin.Context) {
	logger := logger.GetLogger()

	subjects, err := g.repo.FindAll(c)
	if err != nil {
		logger.Error("Failed to fetch subjects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// Create godoc
// @Summary 	create subject
// @Tags 		subjects
// @Accept 		json
// @Produce 	json
// @Param		title 	query 		string 		true 	"Subject_title"
// @Success 	200 	{object} 	object{id=int} 		"OK"
// @Failure 	400 	{object}	models.ApiError 	"Invalid Payload"
// @Failure 	500 	{object} 	models.ApiError
// @Router 		/subjects/{id} [post]
// @Security 	Bearer
func (g *SubjectsHandlers) Create(c *gin.Context) {
	logger := logger.GetLogger()
	var subjectRequest models.Subject

	err := c.BindJSON(&subjectRequest)
	if err != nil {
		logger.Error("Failed JSON binding", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with Json binding"))
		return
	}

	subject := models.Subject{
		Title: subjectRequest.Title,
	}

	id, err := g.repo.Create(c, subject)
	if err != nil {
		logger.Error("Failed to create subject", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})

}

// Update godoc
// @Summary 	update subject
// @Tags 		subjects
// @Accept		json
// @Produce 	json
// @Param 		title 	query 		string 		true 	"Subject_title"
// @Success 	200 	"OK"
// @Failure 	400 	{object} 	models.ApiError 	"Invalid Payload"
// @Failure 	500 	{object} 	models.ApiError
// @Router 		/subjects/{id} [put]
// @Security 	Bearer
func (g *SubjectsHandlers) Update(c *gin.Context) {
	logger := logger.GetLogger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid subject Id format", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Subject not Found"))
		return
	}

	_, err = g.repo.FindById(c, id)
	if err != nil {
		logger.Error("Requested subject not found", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request models.Subject

	err = c.BindJSON(&request)
	if err != nil {
		logger.Error("Failed JSON binding", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with Json binding"))
		return
	}

	updGenre := models.Subject{
		Title: request.Title,
	}

	err = g.repo.Update(c, id, updGenre)
	if err != nil {
		logger.Error("Failed to update subject", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)

}

// Delete godoc
// @Summary 	delete subject by id
// @Tags 		subjects
// @Accept 		json
// @Produce 	json
// @Param 		id 		path 		int 	true 	"Subject_id"
// @Success 	200 	"OK"
// @Failure 	400 	{object} 	models.ApiError "Invalid Payload"
// @Failure 	500 	{object} 	models.ApiError
// @Router 		/subjects/{id} [delete]
// @Security 	Bearer
func (g *SubjectsHandlers) Delete(c *gin.Context) {
	logger := logger.GetLogger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid subject Id format", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid subject ID format"))
	}

	_, err = g.repo.FindById(c, id)
	if err != nil {
		logger.Error("Requested subject not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	err = g.repo.Delete(c, id)
	if err != nil {
		logger.Error("Failed to delete subject", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

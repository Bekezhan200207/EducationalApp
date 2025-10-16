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

type CoursesHandlers struct {
	coursesRepo *repositories.Coursesrepository
}

func NewCoursesHandler(
	coursesRepo *repositories.Coursesrepository) *CoursesHandlers {
	return &CoursesHandlers{
		coursesRepo: coursesRepo,
	}
}

type courseRequest struct {
	Course_title string `json:"title"`
	Description  string `json:"description"`
	Is_published bool   `json:"is_published"`
}

// Create course	godoc
// @Summary 		create course
// @Tags 			courses
// @Accept 			json
// @Produce 		json
// @Param 			title 			query 		string		true	"title"
// @Param 			description 	query 		string 		true 	"description"
// @Param 			is_published 	query 		boolean 	true 	"is_published"
// @Success 		200 			{object} 	object{id=int} 		"OK"
// @Failure 		400 			{object} 	models.ApiError		"error with json dinding"
// @Failure 		500 			{object} 	models.ApiError
// @Router 			/courses [post]
func (h *CoursesHandlers) Create(c *gin.Context) {
	logger := logger.GetLogger()
	var request courseRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("couldn't create course request"))
		logger.Error("Could not bind json data", zap.Error(err))
		return
	}

	course := models.Course{
		Name:         request.Course_title,
		Description:  request.Description,
		Is_published: request.Is_published,
	}

	id, err := h.coursesRepo.Create(c, course)
	if err != nil {
		logger.Error("Failed to create course", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// FindById	godoc
// @Summary 	find course by id
// @Tags 		courses
// @Accept 		json
// @Produce 	json
// @Param 		id 		path 		int 	true 	"id"
// @Success 	200 	{object} 	models.Course "OK"
// @Failure 	400 	{object} 	models.ApiError "invalid course id"
// @Failure 	500 	{object} 	models.ApiError
// @Router 		/courses/{id} [get]
func (h *CoursesHandlers) FindById(c *gin.Context) {
	logger := logger.GetLogger()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid course ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("invalid course id"))
		return
	}

	course, err := h.coursesRepo.FindById(c, id)
	if err != nil {
		logger.Error("Course doesn't exist", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, course)
}

// FindAll godoc
// @Summary 	find all courses
// @Tags 		courses
// @Accept 		json
// @Produce 	json
// @Success 	200 	{object} []models.Course "OK"
// @Failure 	500 	{object} models.ApiError
// @Router 		/courses [get]
func (g *CoursesHandlers) FindAll(c *gin.Context) {
	logger := logger.GetLogger()

	courses, err := g.coursesRepo.FindAll(c)
	if err != nil {
		logger.Error("Failed to fetch courses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, courses)
}

// Update godoc
// @Summary 	update course
// @Tags 		courses
// @Accept 		json
// @Produce 	json
// @Param 		title 			query 		string 		true 	"title"
// @Param 		description 	query 		string 		true 	"description"
// @Param 		is_published	query 		boolean		true 	"is_published"
// @Success 	200  	"OK"
// @Failure 	400 			{object} 	models.ApiError
// @Failure 	500 			{object} 	models.ApiError
// @Router 		/courses/{id} [put]
func (g *CoursesHandlers) Update(c *gin.Context) {
	logger := logger.GetLogger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid course Id format", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Course not Found"))
		return
	}

	_, err = g.coursesRepo.FindById(c, id)
	if err != nil {
		logger.Error("Requested course not found", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request courseRequest

	err = c.BindJSON(&request)
	if err != nil {
		logger.Error("Failed JSON binding", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with Json binding"))
		return
	}

	updCourse := models.Course{
		Name:         request.Course_title,
		Description:  request.Description,
		Is_published: request.Is_published,
	}

	err = g.coursesRepo.Update(c, id, updCourse)
	if err != nil {
		logger.Error("Failed to update course", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)

}

// Delete godoc
// @Summary		delete course by id
// @Tags 		courses
// @Accept 		json
// @Produce 	json
// @Param 		id		path 		int		true 	"id"
// @Success 	200 	"OK"
// @Failure 	400 	{object} 	models.ApiError "Invalid Id"
// @Failure 	500 	{object} 	models.ApiError
// @Router 		/courses/{id} [delete]
func (g *CoursesHandlers) Delete(c *gin.Context) {
	logger := logger.GetLogger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid course Id format", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid course ID format"))
	}

	_, err = g.coursesRepo.FindById(c, id)
	if err != nil {
		logger.Error("Requested course not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	err = g.coursesRepo.Delete(c, id)
	if err != nil {
		logger.Error("Failed to delete course", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

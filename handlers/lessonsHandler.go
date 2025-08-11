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

type LessonsHandler struct {
	lessonsRepo *repositories.Lessonsrepository
}

type lessonRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Subject_id      int    `json:"subject_id"`
	Order           int    `json:"order"`
	Level           string `json:"level"`
	Interest        string `json:"interest"`
	Target_age_min  int    `json:"target_age_min"`
	Target_age_max  int    `json:"target_age_max"`
	Video_data      []byte `json:"video_data"`
	Video_filename  string `json:"video_filename"`
	Video_mime_type string `json:"video_mime_type"`
	Duration_sec    int    `json:"duration"`
	Is_published    bool   `json:"is_published"`
}

func NewLessonsHandler(
	lessonsRepo *repositories.Lessonsrepository) *LessonsHandler {
	return &LessonsHandler{
		lessonsRepo: lessonsRepo,
	}
}

// FindById godoc
// @Summary 	find lesson by id
// @Tags 		lessons
// @Accept 		json
// @Produce 	json
// @Param 		id 		path		int 	true 	"Lesson_id"
// @Success 	200 	{object}	models.Lesson "OK"
// @Failure 	400 	{object}	models.ApiError
// @Failure 	500 	{object}	models.ApiError
// @Router 		/lessons/{id} [get]
func (h *LessonsHandler) FindById(c *gin.Context) {
	logger := logger.GetLogger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid lesson Id", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Lesson Id"))
		return
	}

	lesson, err := h.lessonsRepo.FindById(c, id)

	if err != nil {
		logger.Error("Failed to find lesson", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, lesson)

}

// FindAll godoc
// @Summary 	Find all lessons
// @Tags 		lessons
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} []models.Lesson "OK"
// @Failure 	400 {object} models.ApiError
// @Failure 	500 {object} models.ApiError
// @Router 		/lessons [get]
func (h *LessonsHandler) FindAll(c *gin.Context) {
	logger := logger.GetLogger()

	movies, err := h.lessonsRepo.FindAll(c)
	if err != nil {
		logger.Error("Failed to fetch lessons", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, movies)
}

// Create godoc
// @Summary 	Create Lesson
// @Tags 		lessons
// @Accept	 	json
// @Produce 	json
// @Param		title 				query 		string 			true 	"Lesson_title"
// @Param 		description 		query		string 			true 	"Lesson_description"
// @Param 		subject_id 			query		int 			true 	"Lesson_subject_id"
// @Param 		order 				query		int 			true 	"Topic_order"
// @Param 		level 				query		string 			true 	"Lessons_level" Enum('Beginner', 'Intermediate', 'Advanced')
// @Param 		interest 			query		string 			true 	"Lessons_interest"
// @Param 		target_age_min 		query		int 			true 	"minimum age for viewing lesson"
// @Param 		target_age_max 		query		int 			true 	"maximum age for viewing lesson"
// @Param 		video_data 			query		string 			true 	"Lessons_video"
// @Param 		video_filename 		query		string 			true 	"Lessons_video_filename"
// @Param 		video_mime_type 	query		string 			true 	"Lessons_video_MIME-type"
// @Param 		duration_sec 		query		string 			true 	"Lessons_video_duration_seconds"
// @Param 		is_published 		query		boolean 		true 	"lessons_video_is_published"
// @Success 	200 				{object} 	object{id=int} 	"OK"
// @Failure 	400 				{object}	models.ApiError "Invalid Payload"
// @Failure 	500 				{object} 	models.ApiError
// @Router 		/lessons [post]
func (h *LessonsHandler) Create(c *gin.Context) {
	logger := logger.GetLogger()

	var request lessonRequest

	err := c.Bind(&request)
	if err != nil {
		logger.Error("Failed JSON binding", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with binding"))
		return
	}

	lesson := models.Lesson{
		Title:           request.Title,
		Description:     request.Description,
		Subject_id:      request.Subject_id,
		Order:           request.Order,
		Level:           request.Level,
		Interest:        request.Interest,
		Target_age_min:  request.Target_age_min,
		Target_age_max:  request.Target_age_max,
		Video_data:      request.Video_data,
		Video_filename:  request.Video_filename,
		Video_mime_type: request.Video_mime_type,
		Duration_sec:    request.Duration_sec,
		Is_published:    request.Is_published,
	}

	id, err := h.lessonsRepo.Create(c, lesson)
	if err != nil {
		logger.Error("Failed to create", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Lesson has been created", zap.Int("lesson_id", id))

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// Update godoc
// @Summary Update Lesson
// @Tags lessons
// @Accept json
// @Produce json
// @Param 		id 					path 		int 			true 	"Lesson_id"
// @Param		title 				query 		string 			true 	"Lesson_title"
// @Param 		description 		query		string 			true 	"Lesson_description"
// @Param 		subject_id 			query		int 			true 	"Lesson_subject_id"
// @Param 		order 				query		int 			true 	"Topic_order"
// @Param 		level 				query		string 			true 	"Lessons_level" Enum('Beginner', 'Intermediate', 'Advanced')
// @Param 		interest 			query		string 			true 	"Lessons_interest"
// @Param 		target_age_min 		query		int 			true 	"minimum age for viewing lesson"
// @Param 		target_age_max 		query		int 			true 	"maximum age for viewing lesson"
// @Param 		video_data 			query		string 			true 	"Lessons_video"
// @Param 		video_filename 		query		string 			true 	"Lessons_video_filename"
// @Param 		video_mime_type 	query		string 			true 	"Lessons_video_MIME-type"
// @Param 		duration_sec 		query		string 			true 	"Lessons_video_duration_seconds"
// @Param 		is_published 		query		boolean 		true 	"lessons_video_is_published"
// @Success 	200 				{object} 	object{id=int} 	"OK"
// @Failure 	400 				{object} 	models.ApiError "Invalid Payload"
// @Failure 	500					{object} 	models.ApiError
// @Router 		/lessons/{id} [put]
func (g *LessonsHandler) Update(c *gin.Context) {
	logger := logger.GetLogger()

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		logger.Error("Requested lesson not found", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Lesson not Found"))
		return
	}

	_, err = g.lessonsRepo.FindById(c, id)
	if err != nil {
		logger.Error("Failed to find lesson", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request models.Lesson

	err = c.BindJSON(&request)
	if err != nil {
		logger.Error("failed JSON binding", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with Json binding"))
		return
	}

	updLesson := models.Lesson{
		Title:           request.Title,
		Description:     request.Description,
		Subject_id:      request.Subject_id,
		Order:           request.Order,
		Level:           request.Level,
		Interest:        request.Interest,
		Target_age_min:  request.Target_age_min,
		Target_age_max:  request.Target_age_max,
		Video_data:      request.Video_data,
		Video_filename:  request.Video_filename,
		Video_mime_type: request.Video_mime_type,
		Duration_sec:    request.Duration_sec,
		Is_published:    request.Is_published,
	}

	err = g.lessonsRepo.Update(c, id, updLesson)
	if err != nil {
		logger.Error("Failed to update lesson", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

// Delete godoc
// @Summary 	delete lesson by id
// @Tags 		lessons
// @Accept		json
// @Produce 	json
// @Param 		id 		path 		int 	true 	"Lesson_id"
// @Success 	200 	"OK"
// @Failure 	400 	{object} 	models.ApiError "Invalid Payload"
// @Failure 	500 	{object} 	models.ApiError
// @Router 		/lessons/{id} [delete]
func (g *LessonsHandler) Delete(c *gin.Context) {
	logger := logger.GetLogger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid lesson Id", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid ID"))
		return
	}

	_, err = g.lessonsRepo.FindById(c, id)
	if err != nil {
		logger.Error("Requested lesson not found", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	err = g.lessonsRepo.Delete(c, id)
	if err != nil {
		logger.Error("Failed to delete lesson", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

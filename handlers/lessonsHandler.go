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

type createLessonRequest struct {
	Title           string `form:"title"`
	Description     string `form:"description"`
	Subject_id      int    `form:"subject_id"`
	Order           int    `form:"order"`
	Level           string `form:"level"`
	Interest        string `form:"interest"`
	Target_age_min  int    `form:"target_age_min"`
	Target_age_max  int    `form:"target_age_max"`
	Video_data      []byte `form:"video_data"`
	Video_filename  string `form:"video_filename"`
	Video_mime_type string `form:"video_mime_type"`
	Duration_sec    int    `form:"duration"`
	Is_published    bool   `form:"is_published"`
}

func NewLessonsHandler(
	lessonsRepo *repositories.Lessonsrepository) *LessonsHandler {
	return &LessonsHandler{
		lessonsRepo: lessonsRepo,
	}
}

func (h *LessonsHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	lesson, err := h.lessonsRepo.FindById(c, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
	}

	c.JSON(http.StatusOK, lesson)

}

func (h *LessonsHandler) FindAll(c *gin.Context) {

	movies, err := h.lessonsRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *LessonsHandler) Create(c *gin.Context) {
	var request createLessonRequest

	err := c.Bind(&request)
	if err != nil {
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
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	logger := logger.GetLogger()
	logger.Info("Lesson has been created", zap.Int("movie_id", id))

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (g *LessonsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Lesson not Found"))
	}

	_, err = g.lessonsRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
	}

	var request models.Lesson

	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with Json binding"))
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
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

func (g *LessonsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid ID"))
	}

	_, err = g.lessonsRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
	}

	err = g.lessonsRepo.Delete(c, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

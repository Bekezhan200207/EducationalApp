package handlers

import (
	"go-EdTech/models"
	"go-EdTech/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubjectsHandlers struct {
	repo *repositories.SubjectsRepository
}

func NewSubjectsHandlers(repo *repositories.SubjectsRepository) *SubjectsHandlers {
	return &SubjectsHandlers{repo: repo}
}

func (g *SubjectsHandlers) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Genre not found"))
		return
	}

	subject, err := g.repo.FindById(c, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
	}

	err = c.BindJSON(&subject)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with Json binding"))
	}

	c.JSON(http.StatusOK, subject)

}

func (g *SubjectsHandlers) FindAll(c *gin.Context) {
	genres, err := g.repo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, genres)
}

func (g *SubjectsHandlers) Create(c *gin.Context) {
	var subjectRequest models.Subject

	err := c.BindJSON(&subjectRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with Json binding"))
	}

	subject := models.Subject{
		Title: subjectRequest.Title,
	}

	id, err := g.repo.Create(c, subject)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})

}

func (g *SubjectsHandlers) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Genre not Found"))
	}

	_, err = g.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
	}

	var request models.Subject

	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Error with Json binding"))
	}

	updGenre := models.Subject{
		Title: request.Title,
	}

	err = g.repo.Update(c, id, updGenre)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

func (g *SubjectsHandlers) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid ID"))
	}

	_, err = g.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
	}

	err = g.repo.Delete(c, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

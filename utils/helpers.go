package utils

import (
	"go-EdTech/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


func HandleRepoError(c *gin.Context, logger *zap.Logger, msg string, err error) {
	logger.Error(msg, zap.Error(err))
	c.JSON(http.StatusInternalServerError, models.NewApiError(msg))
}

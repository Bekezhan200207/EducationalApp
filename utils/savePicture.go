package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SavePicture(c *gin.Context, picture *multipart.FileHeader) (string, error) {
	const uploadPath = "images"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(picture.Filename))
	filepath := fmt.Sprintf("%s/%s", uploadPath, filename)

	err := c.SaveUploadedFile(picture, filepath)

	return filename, err
}

package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveUploadedFile(c *gin.Context, file *multipart.FileHeader, targetDir string) (string, error) {
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".pdf" {
		return "", fmt.Errorf("ekstensi file tidak didukung (hanya .jpg, .jpeg, .png, .webp, .pdf)")
	}

	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.ReplaceAll(filepath.Base(file.Filename), " ", "_"), ext)
	filePath := filepath.Join(targetDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return "/" + strings.ReplaceAll(filePath, "\\", "/"), nil
}

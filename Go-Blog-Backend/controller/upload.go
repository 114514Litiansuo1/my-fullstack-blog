package controller

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// protect
	err := c.Request.ParseMultipartForm(5 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 5MB"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}

	// 2. 安全防御：读取文件的真实类型 (Magic Bytes)
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "处理文件失败"})
			return
		}
	}(f)

	buffer := make([]byte, 512) // 读取前 512 字节用于检测类型
	_, err = f.Read(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	contentType := http.DetectContentType(buffer)

	// 只允许真实的图片格式通过
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "非法文件！只允许上传真实图片"})
		return
	}

	// reset
	_, err = f.Seek(0, 0)
	if err != nil {
		return
	}

	// get upload file
	file, err = c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "get upload file failed",
		})
		return
	}

	// check directory of image
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, 0755)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "处理文件失败"})
			return
		}
	}

	// create a unique file name
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), ext)

	// save path
	savePath := filepath.Join(uploadDir, newFileName)

	// save it to storage
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "save image failed",
		})
		return
	}

	// create a url of image
	imageUrl := fmt.Sprintf("http://localhost:8080/uploads/%s", newFileName)

	c.JSON(http.StatusOK, gin.H{
		"message": "upload successfully",
		"url":     imageUrl,
	})
}

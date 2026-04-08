package controller

import (
	"Go-Blog/database"
	"Go-Blog/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	var tags []model.Tag

	if err := database.DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "get Tags failed",
		})
	}

	c.JSON(http.StatusOK, tags)
}

type TagRequest struct {
	Name string `form:"name" binding:"required"`
}

func getTagById(id string) (*model.Tag, error) {
	var tag model.Tag

	if err := database.DB.First(&tag, id).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

func SetTag(c *gin.Context) {
	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter wrong"})
		return
	}

	tag := model.Tag{
		Name: req.Name,
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "set Tag failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "set Tag successfully",
	})
}

func UpdateTag(c *gin.Context) {
	// get Tag id from url
	tagId := c.Param("id")

	// get data from frontend
	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter wrong"})
		return
	}

	// query it if exists
	tag, err := getTagById(tagId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	// update Tag name
	if err := database.DB.Model(&tag).Updates(model.Tag{Name: req.Name}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update Tag failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update Tag successfully",
	})
}

func DeleteTag(c *gin.Context) {
	// get Tag id from url
	tagId := c.Param("id")

	//query it if it exists
	tag, err := getTagById(tagId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	//delete tag
	if err := database.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete Tag failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete Tag successfully",
	})
}

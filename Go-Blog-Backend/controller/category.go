package controller

import (
	"Go-Blog/database"
	"Go-Blog/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []model.Category

	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "get categories failed",
		})
	}

	c.JSON(http.StatusOK, categories)
}

type CategoryRequest struct {
	Name string `form:"name" binding:"required"`
	Desc string `form:"desc"`
}

func getCategoryById(id string) (*model.Category, error) {
	var category model.Category

	if err := database.DB.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func SetCategory(c *gin.Context) {
	var req CategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter wrong"})
		return
	}

	Category := model.Category{
		Name: req.Name,
		Desc: req.Desc,
	}

	if err := database.DB.Create(&Category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "set category failed"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "set category successfully",
	})
}

func UpdateCategory(c *gin.Context) {
	CategoryId := c.Param("id")

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter wrong"})
		return
	}

	category, err := getCategoryById(CategoryId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if err := database.DB.Model(&category).Updates(model.Category{Name: req.Name}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update category failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update category successfully",
	})
}

func DeleteCategory(c *gin.Context) {
	CategoryId := c.Param("id")

	category, err := getCategoryById(CategoryId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete category failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete category successfully",
	})
}

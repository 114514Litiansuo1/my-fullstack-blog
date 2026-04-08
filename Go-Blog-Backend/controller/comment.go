package controller

import (
	"Go-Blog/database"
	"Go-Blog/model"
	"Go-Blog/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostComment(c *gin.Context) {
	ip := c.ClientIP()
	var stats model.IPStats

	// check ip database
	database.DB.FirstOrCreate(&stats, model.IPStats{IP: ip})

	// ban if it was banned
	if stats.IsForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "Your ip is not allowed to comment"})
		return
	}

	// alert if it comment frequently
	if stats.CommentCount >= 5 && time.Since(stats.UpdatedAt) < 1*time.Minute {
		database.DB.Model(&model.IPStats{IP: ip}).Updates(model.IPStats{CommentCount: stats.CommentCount + 1})
		if stats.CommentCount >= 10 {
			database.DB.Model(&stats).Update("IsForbidden", true)
			c.JSON(http.StatusForbidden, gin.H{"error": "Too many request, Your ip is not allowed to comment"})
			return
		}
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many request, you should have a rest"})
		return
	}

	// save comment
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter wrong"})
		return
	}
	comment.IP = ip
	id, err := utils.GenerateId(3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save comment failed"})
		return
	}
	comment.ID = id

	// 开启事务：保存留言并更新 IP 统计
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}

		tx.Model(&stats).Updates(map[string]interface{}{
			"CommentCount": stats.CommentCount + 1,
			"UpdatedAt":    time.Now(),
		})
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save comment failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "save comment successfully",
	})

}

func GetComments(c *gin.Context) {
	var comments []model.Comment
	// 按照创建时间倒序排，最新的留言在最上面
	database.DB.Order("created_at desc").Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}

// delete comment
func DeleteComment(c *gin.Context) {

	var req model.Comment

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "get comment failed",
		})
		return
	}

	if err := database.DB.Where("id = ?", id).Delete(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "delete comment failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete comment successfully",
	})

}

type IPRequest struct {
	IP string `json:"ip"`
}

// block comment ip
func BlockCommentIP(c *gin.Context) {
	var req IPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "parameter wrong",
		})
		return
	}

	var IPStats model.IPStats

	if err := database.DB.FirstOrCreate(&IPStats, model.IPStats{IP: req.IP}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "block failed",
		})
		return
	}

	if err := database.DB.Model(&IPStats).Update("is_forbidden", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "block failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "block successfully",
	})
}

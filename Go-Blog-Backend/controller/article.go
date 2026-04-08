package controller

import (
	"Go-Blog/database"
	"Go-Blog/model"
	"Go-Blog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

func GetSingleArticle(c *gin.Context) {
	var article model.Article

	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "format invalid",
		})
	}

	if err := database.DB.Preload("Category").Preload("Tags").First(&article, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "query failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": article,
	})
}

func GetArticlesWithCursor(c *gin.Context) {
	// query articles with pagination(分页)
	cursorStr := c.DefaultQuery("cursor", "0")
	limitStr := c.DefaultQuery("limit", "10")
	keyword := c.Query("keyword")

	// 原先的指针逻辑
	//cursor, err := strconv.ParseInt(cursorStr, 10, 64)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"message": "transform cursor fail",
	//	})
	//	return
	//}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "transform limit fail",
		})
		return
	}

	var articles []model.Article

	// query database
	query := database.DB.Model(&model.Article{}).Preload("Category").Preload("Tags")
	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "get articles list fail",
		})
		return
	}

	// 模糊搜索
	if keyword != "" {
		searchWord := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ?", searchWord, searchWord)
	}

	if cursorStr != "0" {
		query = query.Where("id < ?", cursorStr)
	}

	query.Order("id desc").Limit(limit).Find(&articles)

	// get nextCursor
	nextCursor := "0"
	if len(articles) > 0 {
		nextCursor = strconv.FormatInt(articles[len(articles)-1].Id, 10)
	}

	// check if it exists more data
	hasMore := len(articles) == limit

	c.JSON(http.StatusOK, gin.H{
		"data":        articles,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}

type ArticleRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary"`
	CategoryId uint   `json:"category_id"`
	TagIds     []uint `json:"tag_ids"`
}

func getArticleById(idStr string) (*model.Article, error) {
	var article model.Article

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := database.DB.First(&article, id).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

func SetArticle(c *gin.Context) {
	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "parameters wrong",
		})
		return
	}

	userId := c.MustGet("userId").(int64)

	// query tags
	var tags []model.Tag

	if len(req.TagIds) > 0 {
		// exist or not
		if err := database.DB.Find(&tags, req.TagIds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "query tags failed",
			})
			return
		}

		// all exist or some
		if len(tags) != len(req.TagIds) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "one or more tags do not exist",
			})
			return
		}
	}

	// create article id
	id, err := utils.GenerateId(2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "create article id failed",
		})
		return
	}

	// categoryid
	if req.CategoryId == 0 {
		req.CategoryId = 1
	}

	// clear script
	p := bluemonday.UGCPolicy()
	req.Content = p.Sanitize(req.Content)

	// create data sending to database
	article := model.Article{
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		CategoryId: req.CategoryId,
		UserId:     userId,
		Tags:       tags,
		Id:         id,
	}
	if err := database.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "create article failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create article successfully",
		"data":    article,
	})
}

func UpdateArticle(c *gin.Context) {
	// get article id from url
	articleId := c.Param("id")

	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "parameters wrong",
		})
		return
	}

	// query article if it exists
	article, err := getArticleById(articleId)

	// update article
	err = database.DB.Model(&article).Updates(model.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CategoryId: req.CategoryId,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "update article failed",
		})
		return
	}

	// tag and article many2many
	if len(req.TagIds) > 0 {
		var tags []model.Tag
		database.DB.Find(&tags, req.TagIds)
		err := database.DB.Model(&article).Association("Tags").Replace(tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update tags failed"})
			return
		}
	} else {
		err := database.DB.Model(&article).Association("Tags").Clear()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "clear article tags failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update article successfully",
	})
}

func DeleteArticle(c *gin.Context) {
	// get article id from url
	articleId := c.Param("id")

	article, err := getArticleById(articleId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "article not found",
		})
		return
	}

	//execute delete
	if err := database.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "delete article failed",
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "delete article successfully",
		})
	}
}

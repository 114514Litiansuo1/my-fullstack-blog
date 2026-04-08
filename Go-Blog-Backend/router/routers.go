package router

import (
	"Go-Blog/controller"
	"Go-Blog/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/api/v1")
	{
		// Alert!!! Delete register API if deploy into the server
		v1.POST("/register", controller.Register)
		v1.POST("/login", middleware.LoginRateLimiter(), controller.Login)

		// Articles Tags and Categories
		v1.GET("/articles", controller.GetArticlesWithCursor)
		v1.GET("/article/:id", controller.GetSingleArticle)
		v1.GET("/tags", controller.GetTags)
		v1.GET("/categories", controller.GetCategories)

		// Post Comment
		v1.POST("/comment", controller.PostComment)
		v1.GET("/comments", controller.GetComments)
	}

	// photo
	r.Static("/uploads", "./uploads")

	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		// Tag API
		auth.POST("/tag", controller.SetTag)
		auth.PUT("/tag/:id", controller.UpdateTag)
		auth.DELETE("/tag/:id", controller.DeleteTag)

		// Category API
		auth.POST("/category", controller.SetCategory)
		auth.PUT("/category/:id", controller.UpdateCategory)
		auth.DELETE("/category/:id", controller.DeleteCategory)

		// Article API
		auth.POST("/article", controller.SetArticle)
		auth.PUT("/article/:id", controller.UpdateArticle)
		auth.DELETE("/article/:id", controller.DeleteArticle)

		// upload files API
		auth.POST("upload", controller.UploadImage)

		// delete comment
		auth.DELETE("/comment/:id", controller.DeleteComment)

		// block ip
		auth.POST("/ip/ban", controller.BlockCommentIP)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "404",
		})
	})

	return r
}

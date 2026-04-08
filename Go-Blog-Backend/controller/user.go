package controller

import (
	"Go-Blog/database"
	"Go-Blog/middleware"
	"Go-Blog/model"
	"Go-Blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRegister struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

func Register(c *gin.Context) {
	var req UserRegister

	//match JSON with UserRegister{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters should contain username, password and email"})
		return
	}

	// query username
	if err := database.DB.Where("username = ?", req.Username).First(&model.User{}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	// use snowflake to create UID
	userId, err := utils.GenerateId(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create id error"})
		return
	}

	// encrypt password
	passwordHash, err := utils.PasswordHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	//save to database
	user := model.User{
		Username: req.Username,
		Password: passwordHash,
		Email:    req.Email,
		Id:       userId,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create user error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successfully"})
}

type UserLogin struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

func Login(c *gin.Context) {
	var req UserLogin

	//match JSON with UserLogin{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters should contain username and password"})
		return
	}

	var user model.User

	//find user in database
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password wrong"})
		return
	}

	//check passwordHash
	if err := utils.CheckPasswordHash(req.Password, user.Password); err != true {
		// 关键：调用限流中间件的记录函数
		middleware.RecordLoginFailure(c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password wrong"})
		return
	}

	//generate jwt if username and password matching with database.
	token, err := utils.GenerateToken(user.Id, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Generate token failed"})
		return
	}

	middleware.ClearLoginFailure(c.ClientIP())

	//send token to frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "login successfully",
		"token":   token,
	})
}

package middleware

import (
	"Go-Blog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Three methods Client to carry a token
		// 1. Request header 2. Request body 3. URI
		// Authorization: Bearer xxxx.xxx.xx
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			c.Abort()
			return
		}

		// split Request body by space
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  http.StatusUnauthorized,
				"error": "Authorization header is invalid",
			})
			c.Abort()
			return
		}

		//Parse(解析) the retrieved(被提取) token string.
		claim, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// save request userID info to context
		c.Set("userId", claim.UserId)
		c.Set("username", claim.Username)

		c.Next()
	}

}

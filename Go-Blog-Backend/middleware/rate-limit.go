package middleware

import (
	"Go-Blog/model"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

func LoginRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		// Redis 里的 Key 格式：login_limit:127.0.0.1
		key := fmt.Sprintf("login_limit: %s", ip)

		// 1. 获取该 IP 失败的次数
		count, err := model.RDB.Get(ctx, key).Int()

		// 如果 Redis 里没有这个 Key，说明该 IP 还没错
		if err != nil && err.Error() != "redis: nil" {
			c.Next()
			fmt.Printf("NoKey IP: %v, Count: %v", key, count)
			return
		}

		// 2. 检查次数是否超过限制（比如 5 次）
		if count >= 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "尝试次数过多，请稍后再试",
			})
			return
		}

		c.Next()
	}
}

// RecordLoginFailure 记录登录失败次数
func RecordLoginFailure(ip string) {
	key := fmt.Sprintf("login_limit: %s", ip)

	// 原子自增 1
	model.RDB.Incr(ctx, key)
	// 设置 5 分钟后过期（即自动解除限制）
	model.RDB.Expire(ctx, key, 1*time.Hour)
}

// ClearLoginFailure 登录成功后清除记录
func ClearLoginFailure(ip string) {
	key := fmt.Sprintf("login_limit:%s", ip)
	model.RDB.Del(ctx, key)
}

package model

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis 地址
		Password: "",           // 密码，没有留空
		DB:       0,            // 默认数据库
	})

	// 测试连接
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic("连接 Redis 失败: " + err.Error())
	}
	fmt.Println("I'm redis")
}

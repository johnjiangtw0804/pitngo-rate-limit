package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/johnjiangtw0804/pitngo-rate-limit/env"
	"github.com/johnjiangtw0804/pitngo-rate-limit/rate_limit"
	"github.com/redis/go-redis/v9"
)

func RegisterRouter(config *env.Configuration, redisClient *redis.Client) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(gin.Logger())   // log 每個請求
	router.Use(gin.Recovery()) // 保護程式不崩潰

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // TODO: I am assuming here we need a load balancer that would sit in front of this RL
		MaxAge:           12 * time.Hour,
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true, // to allow browsers 自帶 credentials
		ExposeHeaders:    []string{"Content-Length"},
	}))

	checkEndpoint := &CheckEndpoint{Limiters: map[string]rate_limit.IRateLimiter{
		"fixed": rate_limit.NewFixedWindowLimiter(redisClient, time.Minute, 10),
		// "token":   tokenLimiter,
		// "sliding": slidingLimiter,
	}}
	router.GET("/api/v1/check", checkEndpoint.CheckHandler())

	return router, nil
}

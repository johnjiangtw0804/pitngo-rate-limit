package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnjiangtw0804/pitngo-rate-limit/rate_limit"
)

type CheckEndpoint struct {
	Limiters map[string]rate_limit.IRateLimiter
}

func (c *CheckEndpoint) CheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Query("userId")
		if userid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
			return
		}

		strategyParam := ctx.Query("strategy")
		rateLimiter, ok := c.Limiters[strategyParam]
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy"})
			return
		}

		isAllowed, err := rateLimiter.Allow(userid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		if !isAllowed {
			// 429 for too many requests
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "allowed"})
	}
}

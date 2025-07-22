package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnjiangtw0804/pitngo-rate-limit/rate_limit"
)

type CheckEndpoint struct {
	Limiter rate_limit.IFixedWindowLimiter
}

func (c *CheckEndpoint) CheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userId")
		if userid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
			return
		}

		isAllowed, err := c.Limiter.Allow(userid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		if !isAllowed {
			// 429
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "allowed"})
	}
}

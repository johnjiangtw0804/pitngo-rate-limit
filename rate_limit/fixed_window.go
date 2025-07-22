package rate_limit

import (
	"fmt"
	"time"

	"github.com/johnjiangtw0804/pitngo-rate-limit/infra"
	"github.com/redis/go-redis/v9"
)

type IFixedWindowLimiter interface {
	Allow(userid string) (bool, error)
}

type FixedWindowLimiter struct {
	RedisOps infra.IFixWindowOps
	Limit    int64
}

func (f *FixedWindowLimiter) Allow(userid string) (bool, error) {
	count, err := f.RedisOps.Set(userid)
	if err != nil {
		return false, fmt.Errorf("Failed Redis Set %w", err)
	}
	if count > f.Limit {
		return false, nil
	}

	return true, nil
}

func NewFixedWindowLimiter(redisClient *redis.Client, WindowSize time.Duration, limit int64) IFixedWindowLimiter {
	fixWindowOps := &infra.FixWindowOps{
		Client:     redisClient,
		WindowSize: WindowSize,
	}
	fix_window_limiter := FixedWindowLimiter{
		RedisOps: fixWindowOps,
		Limit:    limit,
	}

	return &fix_window_limiter
}

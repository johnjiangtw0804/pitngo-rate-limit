package rate_limit

import (
	"fmt"
	"time"

	repository "github.com/johnjiangtw0804/pitngo-rate-limit/repository"
	"github.com/redis/go-redis/v9"
)

type FixedWindowLimiter struct {
	RedisDao repository.IFixWindowRepo
	Limit    int64
}

func (f *FixedWindowLimiter) Allow(userid string) (bool, error) {
	count, err := f.RedisDao.IsAllow(userid)
	if err != nil {
		return false, fmt.Errorf("Failed Redis Set %w", err)
	}
	if count > f.Limit {
		return false, nil
	}

	return true, nil
}

func NewFixedWindowLimiter(redisClient *redis.Client, WindowSize time.Duration, limit int64) IRateLimiter {
	fixWindowLimiter := FixedWindowLimiter{
		RedisDao: &repository.FixWindowRepo{
			Client:     redisClient,
			WindowSize: WindowSize,
		},
		Limit: limit,
	}

	return &fixWindowLimiter
}

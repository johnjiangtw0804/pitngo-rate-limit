package rate_limit

import (
	"fmt"
	"time"

	dao "github.com/johnjiangtw0804/pitngo-rate-limit/dao"
	"github.com/redis/go-redis/v9"
)

type FixedWindowLimiter struct {
	RedisDao dao.IFixWindowDAO
	Limit    int64
}

func (f *FixedWindowLimiter) Allow(userid string) (bool, error) {
	count, err := f.RedisDao.Set(userid)
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
		RedisDao: &dao.FixWindowDao{
			Client:     redisClient,
			WindowSize: WindowSize,
		},
		Limit: limit,
	}

	return &fixWindowLimiter
}

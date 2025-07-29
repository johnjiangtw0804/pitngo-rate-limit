package rate_limit

import (
	"fmt"
	"time"

	"github.com/johnjiangtw0804/pitngo-rate-limit/repository"
	"github.com/redis/go-redis/v9"
)

type SlidingWindowLimiter struct {
	RedisDao repository.ISlidingWindowRepo
	WindowMS int64 // duration in ms
	MaxHits  int64
}

func (s *SlidingWindowLimiter) Allow(userid string) (bool, error) {
	isOK, err := s.RedisDao.IsAllow(userid)
	if err != nil {
		return false, fmt.Errorf("Failed Redis Set %w", err)
	}
	return isOK, nil
}

func NewSlidingWindowLimiter(redisClient *redis.Client, windowDuration time.Duration, maxHits int64) IRateLimiter {
	slidingWindowLimiter := &SlidingWindowLimiter{
		RedisDao: &repository.SlidingWindowRepo{
			Client:  redisClient,
			MaxHits: maxHits,
			// A Duration represents the elapsed time between two instants as an int64 nanosecond count
			WindowMS: int64(windowDuration / time.Millisecond),
		},
		WindowMS: int64(windowDuration / time.Millisecond),
		MaxHits:  maxHits,
	}
	return slidingWindowLimiter
}

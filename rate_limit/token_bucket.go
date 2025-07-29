package rate_limit

import (
	repository "github.com/johnjiangtw0804/pitngo-rate-limit/repository"
	"github.com/redis/go-redis/v9"
)

type TokenBucketLimiter struct {
	RedisDao repository.ITokenBucketRepo
}

func (t *TokenBucketLimiter) Allow(userid string) (bool, error) {
	return t.RedisDao.IsAllow(userid)
}

func NewTokenBucketLimiter(redisClient *redis.Client, refillRate int64, capacity int64) IRateLimiter {
	tokenBucketLimiter := &TokenBucketLimiter{
		RedisDao: &repository.TokenBucketRepo{
			Client:     redisClient,
			RefillRate: refillRate,
			Capacity:   capacity,
		},
	}
	return tokenBucketLimiter
}

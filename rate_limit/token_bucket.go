package rate_limit

import (
	dao "github.com/johnjiangtw0804/pitngo-rate-limit/dao"
	"github.com/redis/go-redis/v9"
)

type TokenBucketLimiter struct {
	RedisDao dao.ITokenBucketDao
}

func (t *TokenBucketLimiter) Allow(userid string) (bool, error) {
	return t.RedisDao.IsAllow(userid)
}

func NewTokenBucketLimiter(redisClient *redis.Client, refillRate int64, capacity int64) IRateLimiter {
	tokenBucketLimiter := &TokenBucketLimiter{
		RedisDao: &dao.TokenBucketDao{
			Client:     redisClient,
			RefillRate: refillRate,
			Capacity:   capacity,
		},
	}
	return tokenBucketLimiter
}

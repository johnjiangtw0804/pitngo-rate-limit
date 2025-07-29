package repository

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ISlidingWindowRepo interface {
	IsAllow(string) (bool, error)
}

type SlidingWindowRepo struct {
	Client   *redis.Client // Make sure to use *redis.Client
	WindowMS int64
	MaxHits  int64
}

func (s *SlidingWindowRepo) IsAllow(userId string) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s:%d:%d:slidingWindow", userId, s.WindowMS, s.MaxHits)
	ctx := context.Background()
	now := time.Now().UnixMilli()
	member := strconv.FormatInt(now, 10) + "-" + strconv.Itoa(rand.Int())

	const luaScript = `
	local key = KEYS[1]
	local window = tonumber(ARGV[1])
	local now = tonumber(ARGV[2])
	local member = ARGV[3]

	redis.call('ZADD', key, now, member)
	redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
	local count = redis.call('ZCARD', key)
	redis.call('PEXPIRE', key, window)
	return count
	`

	result, err := s.Client.Eval(ctx, luaScript, []string{key}, s.WindowMS, now, member).Result()
	if err != nil {
		return false, fmt.Errorf("redis eval error: %w", err)
	}

	// Result can be int64 or float64 depending on Redis client behavior
	var count int64
	switch v := result.(type) {
	case int64:
		count = v
	case float64:
		count = int64(v)
	default:
		return false, fmt.Errorf("unexpected result type: %T", result)
	}

	return count <= s.MaxHits, nil
}

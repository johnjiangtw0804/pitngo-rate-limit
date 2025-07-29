package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ITokenBucketRepo interface {
	IsAllow(userId string) (bool, error)
}

type TokenBucketRepo struct {
	Client     *redis.Client
	RefillRate int64 // in sec
	Capacity   int64
}

func (t *TokenBucketRepo) IsAllow(userId string) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s:tokenBucket:%d:%d", userId, t.RefillRate, t.Capacity)
	// Redis structure
	//
	//	token_bucket:{user_id} -> {
	//	  "tokens": 8,
	//	  "last_refill_timestamp": 1721905000
	//	}
	const luaScript = `
	local key = KEYS[1]
	local now = tonumber(ARGV[1])
	local refill_rate_in_sec = tonumber(ARGV[2])
	local capacity = tonumber(ARGV[3]) or 100

	local data = redis.call("HMGET", key, "tokens", "last_refill_timestamp")
	local current_tokens = tonumber(data[1]) or 0
	local last_refill_timestamp = tonumber(data[2]) or now

	local refill_tokens = math.floor((now - last_refill_timestamp) * refill_rate_in_sec)
	local newToken = math.min(capacity, refill_tokens + current_tokens)

	local is_allow = 0

	if newToken > 0 then
		newToken = newToken - 1
		is_allow = 1
	end

	redis.call("HSET", key, "tokens", newToken, "last_refill_timestamp", now)
	return is_allow
	`
	ctx := context.Background()
	now := strconv.FormatInt(time.Now().Unix(), 10)
	refill := strconv.FormatInt(t.RefillRate, 10)
	capacity := strconv.FormatInt(t.Capacity, 10)
	result, err := t.Client.Eval(ctx, luaScript, []string{key}, []string{now, refill, capacity}).Result()
	if err != nil {
		return false, err
	}

	allowed, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected Redis return type: %T", result)
	}

	return allowed == 1, nil
}

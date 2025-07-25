package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type IFixWindowDAO interface {
	// we need to set the rate_limit:userid:currentTimeToMins
	Set(string) (int64, error)
}

type FixWindowDao struct {
	Client     *redis.Client
	WindowSize time.Duration
}

func (f *FixWindowDao) Set(userId string) (int64, error) {
	key := fmt.Sprintf("rate_limit:fixedWindow:%s:%s", userId, time.Now().Format("200601021504"))

	const luaScript = `
	local counter = redis.call("INCR", KEYS[1])
	if tonumber(counter) == 1 then
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end
	return counter
	`
	ctx := context.Background()
	result, err := f.Client.Eval(ctx, luaScript, []string{key}, strconv.Itoa(int(f.WindowSize.Seconds()))).Result()
	if err != nil {
		return 0, err
	}

	count, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected result type: %T", result)
	}
	return count, nil
}

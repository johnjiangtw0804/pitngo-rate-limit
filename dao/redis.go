package dao

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/johnjiangtw0804/pitngo-rate-limit/env"
	"github.com/redis/go-redis/v9"
)

func ConnectDBS(config *env.Configuration) (*redis.Client, error) {
	dbID, err := strconv.Atoi(config.RedisDB)
	if err != nil {
		return nil, fmt.Errorf("invalid RedisDB: %v", err)
	}

	url := fmt.Sprintf("redis://%s:%s/%d", config.RedisHost, config.RedisPort, dbID)
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %v", err)
	}

	log.Println("Redis connected: OK")
	return client, nil
}

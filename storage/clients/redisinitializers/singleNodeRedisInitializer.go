package redisinitializers

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v7"
)

type SingleNodeRedisInitializer struct {
}

func (s SingleNodeRedisInitializer) Initialize(ctx context.Context, options interface{})  (*redis.Client, error)  {
	redisOption, ok := options.(redis.Options)
	if !ok {
		return nil, errors.New("Can't cast options to redis Options")
	}
	if redisOption.Addr == "" {
		return nil, errors.New("Addr can't be empty")
	}
	client := redis.NewClient(&redisOption)
	_, err := client.Ping().Result()
	return client, err
}

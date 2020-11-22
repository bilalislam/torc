package redisinitializers

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v7"
)

type FailoverRedisInitializer struct {
}

func (f FailoverRedisInitializer) Initialize(ctx context.Context, options interface{}) (*redis.Client, error) {
	redisOption, ok := options.(redis.FailoverOptions)
	if !ok {
		return nil, errors.New("Can't cast options to redis FailoverOptions")
	}
	if len(redisOption.SentinelAddrs) == 0 {
		return nil, errors.New("SentinelAddrs can not be empty")
	}

	if redisOption.MasterName == "" {
		panic("redis.sentinelMasterName can not be empty!")
	}
	client := redis.NewFailoverClient(&redisOption)
	_, err := client.Ping().Result()
	return client, err
}

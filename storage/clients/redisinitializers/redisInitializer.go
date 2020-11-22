package redisinitializers

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
)

type ConnectType string

const (
	Failover   ConnectType = "failover"
	SingleNode ConnectType = "singleNode"
)

type RedisInitializer interface {
	Initialize(ctx context.Context, options interface{}) (*redis.Client, error)
}

var initializers = map[string]RedisInitializer{
	"failover":   FailoverRedisInitializer{},
	"singleNode": SingleNodeRedisInitializer{},
}

func GetInitializer(connectType ConnectType) RedisInitializer {
	initializer, exist := initializers[string(connectType)]
	if !exist {
		panic(fmt.Sprintf("Initializer not found with selection! %v", connectType))
	}
	return initializer
}

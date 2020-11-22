package clients

import (
	"context"
	"encoding/json"
	"github.com/bilalislam/torc/storage"
	"github.com/bilalislam/torc/storage/clients/redisinitializers"
	"github.com/bilalislam/torc/storage/models"
	"github.com/go-redis/redis/v7"
)

type IRedisRepository interface {
	storage.IRepository
	SaveAsUnstructured(id string, data []byte) error
	SetKeyPrefix(key string)
	GetKeyPrefix() string
	GetByIdAsString(id string) (string, error)
}
type RedisRepository struct {
	Conn      *redis.Client
	KeyPrefix string
}

func (r *RedisRepository) Connect(ctx context.Context, dbContext interface{}) error {
	initializer := redisinitializers.GetInitializer(redisinitializers.SingleNode)
	client, err := initializer.Initialize(ctx, dbContext)
	r.Conn = client
	return err
}

func (r *RedisRepository) ConnectFailoverCluster(ctx context.Context, dbContext interface{}) error {
	initializer := redisinitializers.GetInitializer(redisinitializers.Failover)
	client, err := initializer.Initialize(ctx, dbContext)
	r.Conn = client
	return err
}

func (r *RedisRepository) SetKeyPrefix(key string) {
	r.KeyPrefix = key
}

func (r *RedisRepository) GetKeyPrefix() string {
	return r.KeyPrefix
}

func (r *RedisRepository) createKey(id string) string {
	return r.KeyPrefix + id
}

func (r *RedisRepository) GetByIdAsString(id string) (string, error) {
	result, err := r.Conn.Get(r.createKey(id)).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}


func (r *RedisRepository) SaveAsUnstructured(id string, data []byte) error {
	err := r.Conn.Set(r.createKey(id), data, 0).Err()
	return err
}

func (r *RedisRepository) GetById(ctx context.Context, id string, model models.IModel) error {
	result, err := r.Conn.Get(r.createKey(id)).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), model)
}

func (r *RedisRepository) Save(ctx context.Context, model models.IModel) error {
	serialized, err := json.Marshal(model)
	if err != nil {
		return err
	}
	err = r.Conn.Set(r.createKey(model.GetId()), serialized, 0).Err()
	return err
}

func (r *RedisRepository) Update(ctx context.Context, model models.IModel) (int64, error) {
	serialized, err := json.Marshal(model)
	if err != nil {
		return 0, err
	}
	err = r.Conn.Set(r.createKey(model.GetId()), serialized, 0).Err()
	return 1, err
}

func (r *RedisRepository) Delete(ctx context.Context, id string) (int64, error) {
	affectedCount, err := r.Conn.Del(r.createKey(id)).Result()
	if err != nil {
		return 0, err
	}
	return affectedCount, nil
}

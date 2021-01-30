package cache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/nicuf/file-processor-api/task"
)

type redisCache struct {
	redisClient *redis.Client
}

type Cache interface {
	Set(key string, value task.Task) error
	Get(key string) (*task.Task, error)
	GetNextID() (string, error)
}

func NewRedisCache(rc *redis.Client) Cache {
	return &redisCache{rc}
}

func (r *redisCache) Set(key string, value task.Task) error {
	json, err := value.ToJSON()
	if err != nil {
		return err
	}
	err = r.redisClient.Set(key, string(json), 0).Err()
	return err
}

func (r *redisCache) Get(key string) (*task.Task, error) {
	jsonString, err := r.redisClient.Get(key).Result()
	if err != nil {
		return nil, err
	}
	val, err := task.FromJson(jsonString)
	if err != nil {
		return nil, err
	}
	return val, err
}

func (r *redisCache) GetNextID() (string, error) {
	nextID, err := r.redisClient.Incr("nextID").Result()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", nextID), nil
}

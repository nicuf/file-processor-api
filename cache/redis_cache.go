package cache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/nicuf/file-processor-api/task"
)

type redisCache struct {
	addr         string        `default:"localhost:6379"`
	password     string        `default:""`
	db           int           `default:0`
	queueChannel string        `default:"FileProcessor"`
	redisClient  *redis.Client `default: nil`
}

type Cache interface {
	Set(key string, value task.Task) error
	Get(key string) (*task.Task, error)
	GetNextID() (string, error)
	Subscribe() (<-chan *redis.Message, error)
	Publish(taskUUID string) error
}

func NewRedisCache() Cache {
	cache := redisCache{}
	cache.redisClient = redis.NewClient(&redis.Options{
		Addr:     cache.addr,
		Password: cache.password,
		DB:       cache.db,
	})

	return &cache
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

func (r *redisCache) Subscribe() (<-chan *redis.Message, error) {
	pubsub := r.redisClient.Subscribe(r.queueChannel)
	_, err := pubsub.Receive()
	if err != nil {
		return nil, err
	}
	return pubsub.Channel(), err
}

func (r *redisCache) Publish(taskUUID string) error {
	err := r.redisClient.Publish(r.queueChannel, taskUUID).Err()
	if err != nil {
		return err
	}
	return err
}

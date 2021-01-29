package cache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/nicuf/file-processor-api/task"
)

type redisCache struct {
	addr                     string        `default:"localhost:6379"`
	password                 string        `default:""`
	db                       int           `default:0`
	apiMessagesChanName      string        `default:"ApiTaskMessages"`
	availableWorkersChanName string        `default:"AvailableWorkers"`
	redisClient              *redis.Client `default: nil`
}

type Cache interface {
	Set(key string, value task.Task) error
	Get(key string) (*task.Task, error)
	GetNextID() (string, error)

	SubscribeToApiMessages() (<-chan *redis.Message, error)
	PublishTaskMessage(taskUUID string) error

	SubscribeToAvailableWorkers() (<-chan *redis.Message, error)
	PublishAvailableWorker(workerUUID string) error

	SubscribeWorker(workerUUID string) (<-chan *redis.Message, error)
	PublishWork(workerUUID string, taskUUID string) error
}

func NewRedisCache() Cache {
	cache := redisCache{
		addr:                     "localhost:6379",
		password:                 "",
		db:                       0,
		apiMessagesChanName:      "ApiTaskMessages",
		availableWorkersChanName: "AvailableWorkers",
	}
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

func (r *redisCache) SubscribeToApiMessages() (<-chan *redis.Message, error) {
	pubsub := r.redisClient.Subscribe(r.apiMessagesChanName)
	return pubsub.Channel(), nil
}

func (r *redisCache) PublishTaskMessage(taskUUID string) error {
	err := r.redisClient.Publish(r.apiMessagesChanName, taskUUID).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *redisCache) SubscribeToAvailableWorkers() (<-chan *redis.Message, error) {
	pubsub := r.redisClient.Subscribe(r.availableWorkersChanName)
	return pubsub.Channel(), nil
}

func (r *redisCache) PublishAvailableWorker(workerUUID string) error {
	err := r.redisClient.Publish(r.availableWorkersChanName, workerUUID).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *redisCache) SubscribeWorker(workerUUID string) (<-chan *redis.Message, error) {
	pubsub := r.redisClient.Subscribe(workerUUID)
	return pubsub.Channel(), nil
}

func (r *redisCache) PublishWork(workerUUID string, taskUUID string) error {
	err := r.redisClient.Publish(workerUUID, taskUUID).Err()
	if err != nil {
		return err
	}
	return err
}

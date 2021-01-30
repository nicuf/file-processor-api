package message_queue

import "github.com/go-redis/redis"

type redisMessageQueue struct {
	apiMessagesChanName      string
	availableWorkersChanName string `default:"AvailableWorkers"`
	redisClient              *redis.Client
}

func NewRedisMessageQueue(rc *redis.Client) MessageQueue {
	return &redisMessageQueue{
		apiMessagesChanName:      "ApiTaskMessages",
		availableWorkersChanName: "AvailableWorkers",
		redisClient:              rc,
	}
}

type MessageQueue interface {
	SubscribeToApiMessages() (<-chan *redis.Message, error)
	PublishTaskMessage(taskUUID string) error

	SubscribeToAvailableWorkers() (<-chan *redis.Message, error)
	PublishAvailableWorker(workerUUID string) error

	SubscribeWorker(workerUUID string) (<-chan *redis.Message, error)
	PublishWork(workerUUID string, taskUUID string) error
}

func (r *redisMessageQueue) SubscribeToApiMessages() (<-chan *redis.Message, error) {
	pubsub := r.redisClient.Subscribe(r.apiMessagesChanName)
	return pubsub.Channel(), nil
}

func (r *redisMessageQueue) PublishTaskMessage(taskUUID string) error {
	err := r.redisClient.Publish(r.apiMessagesChanName, taskUUID).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *redisMessageQueue) SubscribeToAvailableWorkers() (<-chan *redis.Message, error) {
	pubsub := r.redisClient.Subscribe(r.availableWorkersChanName)
	return pubsub.Channel(), nil
}

func (r *redisMessageQueue) PublishAvailableWorker(workerUUID string) error {
	err := r.redisClient.Publish(r.availableWorkersChanName, workerUUID).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *redisMessageQueue) SubscribeWorker(workerUUID string) (<-chan *redis.Message, error) {
	pubsub := r.redisClient.Subscribe(workerUUID)
	return pubsub.Channel(), nil
}

func (r *redisMessageQueue) PublishWork(workerUUID string, taskUUID string) error {
	err := r.redisClient.Publish(workerUUID, taskUUID).Err()
	if err != nil {
		return err
	}
	return err
}

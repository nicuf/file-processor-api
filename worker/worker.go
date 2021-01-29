package worker

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nicuf/file-processor-api/cache"
)

type workerPool struct {
	workersNumber    int
	messageQueue     cache.Cache
	log              *log.Logger
	availableWorkers chan string
}

type WorkerPool interface {
	StartWorkers()
	StartMaster()
}

func NewWorkerPool(numberOfWorkers int, messageQueue cache.Cache, log *log.Logger) WorkerPool {
	return &workerPool{
		workersNumber: numberOfWorkers,
		messageQueue:  messageQueue,
		log:           log,
	}
}

func (wp *workerPool) StartWorkers() {
	for i := 0; i < wp.workersNumber; i++ {
		go func() {

			uuid := uuid.New().String()
			wp.log.Println("Starting worker: ", uuid)
			err := wp.messageQueue.PublishAvailableWorker(uuid)
			if err != nil {
				wp.log.Fatal("Unable to publish available worker, ", err)
			}

			ch, err := wp.messageQueue.SubscribeWorker(uuid)
			if err != nil {
				log.Fatal("Unable to subscribe to the worker queue ", err)
			}

			for {
				wp.log.Println("Waiting for work, worker: ", uuid)
				message := <-ch
				wp.log.Printf("Worker: %v received task: %v", uuid, message.Payload)
				wp.performTask(message.Payload)
				err = wp.messageQueue.PublishAvailableWorker(uuid)
				if err != nil {
					wp.log.Fatal("Unable to publish available worker, ", err)
				}
			}
		}()
	}
}

func (wp *workerPool) StartMaster() {
	wp.availableWorkers = make(chan string, wp.workersNumber)
	go func() {
		availableWorkersChan, err := wp.messageQueue.SubscribeToAvailableWorkers()
		if err != nil {
			log.Fatal("Unable to subscribe to available workers, ", err)
		}
		for {
			availableWorker := <-availableWorkersChan
			wp.log.Println("Received available worker: ", availableWorker.Payload)
			wp.availableWorkers <- availableWorker.Payload
		}
	}()

	go func() {
		wp.log.Println("Starting Master")
		messageChan, err := wp.messageQueue.SubscribeToApiMessages()
		if err != nil {
			log.Fatal("Unable to subscribe to api message queue, ", err)
		}

		for {
			fileUUID := <-messageChan
			wp.log.Println("Received api message: ", fileUUID.Payload)
			availableWorker := <-wp.availableWorkers

			wp.log.Printf("Assigning taskUUID: %v to worker: %v\n", fileUUID.Payload, availableWorker)
			err = wp.messageQueue.PublishWork(availableWorker, fileUUID.Payload)

			if err != nil {

				wp.log.Printf("Unable to assign taskUUID: %v to worker: %v, %v \n", fileUUID.Payload, availableWorker, err)
			}

		}
	}()
}

func (wp *workerPool) performTask(uuid string) error {
	wp.log.Println("Processing: ", uuid)
	time.Sleep(50 * time.Second)
	wp.log.Println("Finished: ", uuid)
	return nil
}

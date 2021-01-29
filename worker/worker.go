package worker

import (
	"log"

	"github.com/google/uuid"
	"github.com/nicuf/file-processor-api/cache"
)

type workerPool struct {
	workersNumber    int
	messageQueue     cache.Cache
	log              *log.Logger
	work             func(fileUUID string) error
	availableWorkers chan string
}

type WorkerPool interface {
	StartWorkers()
	StartMaster()
}

func NewWorkerPool(numberOfWorkers int, messageQueue cache.Cache, log *log.Logger, work func(fileUUID string) error) WorkerPool {
	return &workerPool{
		workersNumber: numberOfWorkers,
		messageQueue:  messageQueue,
		log:           log,
		work:          work,
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
				err = wp.performTask(message.Payload)
				if err != nil {
					wp.log.Printf("Worker: %v failed task: %v", uuid, message.Payload)
				} else {
					wp.log.Printf("Worker: %v finished task: %v", uuid, message.Payload)
				}
				err = wp.messageQueue.PublishAvailableWorker(uuid)
				if err != nil {
					wp.log.Fatal("Unable to publish available worker, ", err)
				}
			}
		}()
	}
}

func (wp *workerPool) StartMaster() {
	wp.log.Println("Starting Master")
	initChan := make(chan int, 2)
	wp.availableWorkers = make(chan string, wp.workersNumber)
	go func() {
		availableWorkersChan, err := wp.messageQueue.SubscribeToAvailableWorkers()
		if err != nil {
			log.Fatal("Unable to subscribe to available workers, ", err)
		}
		initChan <- 1
		for {
			availableWorker := <-availableWorkersChan
			wp.log.Println("Received available worker: ", availableWorker.Payload)
			wp.availableWorkers <- availableWorker.Payload
		}
	}()

	go func() {
		messageChan, err := wp.messageQueue.SubscribeToApiMessages()
		if err != nil {
			log.Fatal("Unable to subscribe to api message queue, ", err)
		}
		initChan <- 2
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

	<-initChan
	<-initChan
}

func (wp *workerPool) performTask(uuid string) error {
	err := wp.work(uuid)
	if err != nil {
		return err
	}
	return nil
}

package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/nicuf/file-processor-api/cache"
	"github.com/nicuf/file-processor-api/message_queue"
	"github.com/nicuf/file-processor-api/task"
	"github.com/nicuf/file-processor-api/worker"
)

type TaskHandler struct {
	log          *log.Logger
	cache        cache.Cache
	messageQueue message_queue.MessageQueue
}

func NewTaskHandler(l *log.Logger, c cache.Cache, m message_queue.MessageQueue) *TaskHandler {
	return &TaskHandler{l, c, m}
}

func (taskHandler *TaskHandler) addNewTask(uuid string) (*task.Task, error) {
	id, err := taskHandler.cache.GetNextID()
	if err != nil {
		return nil, err
	}
	currentTask := task.Task{
		ID:           id,
		CreationTime: time.Now().Format(time.RFC3339),
		FileUUID:     uuid,
		Status:       task.Queued,
		Result:       []string{},
	}

	err = taskHandler.cache.Set(uuid, currentTask)
	if err != nil {
		return nil, err
	}

	err = taskHandler.messageQueue.PublishTaskMessage(uuid)
	if err != nil {
		return nil, err
	}

	return &currentTask, nil
}

//Send a file to analysis that is specified via the ID
// swagger:route POST /task/{uuid} file-processor task
// Add a task to be processed
// responses:
//	200: taskInfoResponse
func (taskHandler *TaskHandler) AddTask(rw http.ResponseWriter, req *http.Request) {
	taskHandler.log.Println("Handle POST Task.")
	rw.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(req)
	uuid := vars["uuid"]

	currentTask, err := taskHandler.cache.Get(uuid)
	if err == redis.Nil {
		currentTask, err = taskHandler.addNewTask(uuid)
	}
	if err != nil {
		http.Error(rw, "Unable to add new task.", http.StatusInternalServerError)
		taskHandler.log.Println("Error adding new task: ", err)
	}
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(currentTask)
	if err != nil {
		http.Error(rw, "Unable to add new task.", http.StatusInternalServerError)
		taskHandler.log.Println("Error adding new task: ", err)
	}
}

//Get information about an analysis (execution status and results)
// swagger:route GET /task/{uuid} file-processor task
// Returns task info
// responses:
//  200: taskInfoResponse
func (taskHandler *TaskHandler) GetTaskInfo(rw http.ResponseWriter, req *http.Request) {
	taskHandler.log.Println("Handle Get Task.")
	rw.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(req)
	uuid := vars["uuid"]

	currentTask, err := taskHandler.cache.Get(uuid)
	if err == redis.Nil {
		http.Error(rw, "File UUID does not exist in the task queue.", http.StatusNotFound)
	}
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(currentTask)
	if err != nil {
		http.Error(rw, "Unable to process request.", http.StatusInternalServerError)
	}

}

//Search for the files that contain a particular UUID
// swagger:route GET /files/{uuid} file-processor  searchFiles
// Returns a list of files with UUID
// responses:
//  200: searchFilesResponse
func (taskHandler *TaskHandler) SearchFiles(rw http.ResponseWriter, req *http.Request) {

}

// Given a file UUID, check if there’s a chain of referenced files that starts with that file and contains a loop
// swagger:route GET /isloop/{uuid} file-processor isLoop
// Returns true or false
// responses:
//  200: isLoopResponse
func (taskHandler *TaskHandler) IsLoop(rw http.ResponseWriter, req *http.Request) {
	taskHandler.log.Println("Handle Get IsLoop.")
	rw.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(req)
	uuid := vars["uuid"]

	val, err := worker.IsLoop(uuid)
	if err != nil {
		http.Error(rw, "Unable to process request.", http.StatusInternalServerError)
	}
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(val)
	if err != nil {
		http.Error(rw, "Unable to process request.", http.StatusInternalServerError)
	}
}

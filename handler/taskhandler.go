package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nicuf/file-processor-api/task"
)

type TaskHandler struct {
	log *log.Logger
}

func NewTaskHandler(l *log.Logger) *TaskHandler {
	return &TaskHandler{l}
}

//Send a file to analysis that is specified via the ID
// swagger:route POST /addtask/{uuid} addtask addTask
// Add a task to be processed
// responses:
//	200: taskInfoResponse
func (taskHandler *TaskHandler) AddTask(rw http.ResponseWriter, req *http.Request) {
	taskHandler.log.Println("Handle POST Task.")
	rw.Header().Add("Content-Type", "application/json")
	status := task.Queued
	encoder := json.NewEncoder(rw)
	err := encoder.Encode(status)
	if err != nil {
		http.Error(rw, "Unable to encode response", http.StatusInternalServerError)
	}
}

//Get information about an analysis (execution status and results)
// swagger:route GET /task/{id} task getTaskInfo
// Returns task info
// responses:
//  200: taskInfoResponse
func (taskHandler *TaskHandler) GetTaskInfo(rw http.ResponseWriter, req *http.Request) {
	taskHandler.log.Println("Handle Get Task.")
	rw.Header().Add("Content-Type", "application/json")
	task := task.Task{
		ID:           "1",
		CreationTime: "today",
		FileUUID:     "123e4567-e89b-12d3-a456-426655440000",
		Status:       task.Finished,
		Result:       []string{"first", "second", "third"},
	}
	encoder := json.NewEncoder(rw)
	err := encoder.Encode(task)
	if err != nil {
		http.Error(rw, "Unable to encode response", http.StatusInternalServerError)
	}
}

//Search for the files that contain a particular UUID
// swagger:route GET /files/{uuid} files searchFiles
// Returns a list of files with UUID
// responses:
//  200: searchFilesResponse
func (taskHandler *TaskHandler) SearchFiles(rw http.ResponseWriter, req *http.Request) {

}

package handler

import (
	"log"
	"net/http"
)

type TaskHandler struct {
	log *log.Logger
}

func NewTaskHandler(l *log.Logger) *TaskHandler {
	return &TaskHandler{l}
}

func (taskHandler *TaskHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		taskHandler.addTask(rw, req)
		return
	}

	if req.Method == http.MethodGet {
		taskHandler.getTaskInfo(rw, req)
		return
	}
}

//Send a file to analysis that is specified via the ID
// swagger:route POST /addtask/{uuid} addtask addTask
// Add a task to be processed
// responses:
//  default: genericError
//	200: taskStatusResponse
func (taskHandler *TaskHandler) addTask(rw http.ResponseWriter, req *http.Request) {

}

//Get information about an analysis (execution status and results)
// swagger:route GET /task/{id} task getTaskInfo
// Returns task info
// responses:
//  default: genericError
//  200: taskInfoResponse
func (taskHandler *TaskHandler) getTaskInfo(rw http.ResponseWriter, req *http.Request) {

}

//Search for the files that contain a particular UUID
// swagger:route GET /files/{uuid} files searchFiles
// Returns a list of files with UUID
// responses:
//  default: genericError
//  200: searchFilesResponse
func (taskHandler *TaskHandler) searchFiles(rw http.ResponseWriter, req *http.Request) {

}

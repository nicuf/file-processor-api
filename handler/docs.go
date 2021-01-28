//	Package handler Classification of File Processor API
//
//	Documentation for File Processor API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//	swagger:meta
package handler

import "github.com/nicuf/file-processor-api/task"

// swagger:response taskStatusResponse
type taskStatusResponseWrapper struct {
	// Status of the task
	// in: body
	Body task.TaskStatus
}

// swagger:response taskInfoResponse
type taskInfoResponseWrapper struct {
	// Task info
	// in: body
	Body task.Task
}

// swagger:response searchFilesResponse
type searchFilesResponseWrapper struct {
	// The list of file that contains specific UUID
	// in: body
	Body []string `json:"files"`
}

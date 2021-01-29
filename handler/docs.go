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

// Task info
// swagger:response taskInfoResponse
type taskInfoResponseWrapper struct {
	// Task info
	// in: body
	Body task.Task `json:"task"`
}

// Array of files that contains a specific UUID
// swagger:response searchFilesResponse
type searchFilesResponseWrapper struct {
	// The list of file that contains specific UUID
	// in: body
	Body []string `json:"files"`
}

// True or False
// swagger:response isLoopResponse
type isLoopResponseWrapper struct {
	// in: body
	Body bool `json:"result"`
}

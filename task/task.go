package task

//TaskStatus is the status of a Task
//swagger:model
type TaskStatus int

const (
	Started TaskStatus = iota
	Finished
	Queued
	Failed
)

// Task defines the structure for an API Task
// swagger:model
type Task struct {
	// The id of the task
	// required: false
	// min: 1
	ID string `json:"id"`

	// The creation time of the task
	// required: false
	// max length: 100
	CreationTime string `json:"creation_time"`

	// The UUID of the file that is requested to be processed
	// required: true
	// max length: 100
	FileUUID string `json:"file_uuid"`

	// The status of the task
	// required: false
	// min: 0
	Status TaskStatus `json:"status"`

	// The results contains the UUIDs that are in the file that was processed
	//required: false
	Result []string `json:"result"`
}

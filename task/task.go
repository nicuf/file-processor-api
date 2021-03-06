package task

import "encoding/json"

//TaskStatus is the status of a Task
//swagger:model
type TaskStatus string

const (
	Queued   TaskStatus = "Queued"
	Started  TaskStatus = "Started"
	Finished TaskStatus = "Finished"
	Failed   TaskStatus = "Failed"
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
	// max length: 10
	Status TaskStatus `json:"status"`

	// The results contains the UUIDs that are in the file that was processed
	//required: false
	Result []string `json:"result"`
}

func (t *Task) ToJSON() (string, error) {
	jsonString, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(jsonString), err
}

func FromJson(jsonString string) (*Task, error) {
	val := Task{}
	err := json.Unmarshal([]byte(jsonString), &val)
	if err != nil {
		return nil, err
	}
	return &val, err
}

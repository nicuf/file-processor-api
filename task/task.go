package task

//TaskStatus is the status of a Task
type TaskStatus int

const (
	Started TaskStatus = iota
	Finished
	Queued
	Failed
)

type Task struct {
	ID           string
	CreationDate string
	FileID       string
	Status       TaskStatus
	Result       []string
}

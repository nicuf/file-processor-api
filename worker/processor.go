package worker

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"github.com/nicuf/file-processor-api/cache"
	"github.com/nicuf/file-processor-api/task"
)

type processor struct {
	log   *log.Logger
	cache cache.Cache
}

type Processor interface {
	RunTask(fileUUID string) error
}

func NewProcessor(l *log.Logger, c cache.Cache) Processor {
	return &processor{l, c}
}

func (p *processor) RunTask(fileUUID string) error {

	p.updateCache(fileUUID, task.Started, []string{})

	path := filepath.Join("input", fileUUID)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		p.log.Println("File reading error", err)
		p.updateCache(fileUUID, task.Failed, []string{})
		return err
	}

	re := regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")

	uuids := re.FindAllString(string(data), -1)
	p.log.Println("Found ids: ", uuids)
	return p.updateCache(fileUUID, task.Finished, uuids)
}

func (p *processor) updateCache(fileUUID string, status task.TaskStatus, result []string) error {

	currentTask, err := p.cache.Get(fileUUID)
	if err != nil {
		p.log.Println("Unable to retreive from cache the task with uuid:", fileUUID)
		return err
	}

	currentTask.Status = status

	currentTask.Result = append(currentTask.Result, result...)

	err = p.cache.Set(fileUUID, *currentTask)
	if err != nil {
		p.log.Println("Unable to put in cache the task with uuid: ", fileUUID)
		return err
	}
	return nil
}

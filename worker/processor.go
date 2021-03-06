package worker

import (
	"io/ioutil"
	"log"
	"os"
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
	uuids, err := readUUIDSFromFile(fileUUID)
	if err != nil {
		p.log.Println("File reading error", err)
		p.updateCache(fileUUID, task.Failed, []string{})
		return err
	}
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

func IsLoop(fileUUID string) (bool, error) {
	fileGraph, err := getFileGraph(fileUUID)
	if err != nil {
		return false, err
	}
	visited := make(map[string]bool)
	return isLoop(fileGraph, visited, fileUUID), nil
}

func isLoop(fileGraph map[string][]string, visited map[string]bool, currentNode string) bool {
	v, ok := visited[currentNode]
	if ok && v {
		return true
	}

	adjacentNodes, _ := fileGraph[currentNode]
	visited[currentNode] = true
	for _, v := range adjacentNodes {
		result := isLoop(fileGraph, visited, v)
		if result == true {
			return true
		}
	}

	visited[currentNode] = false
	return false
}

func getFileGraph(fileUUID string) (map[string][]string, error) {
	fileGraph := make(map[string][]string)
	filesToRead := []string{fileUUID}

	for len(filesToRead) > 0 {
		currentFile := filesToRead[0]
		filesToRead = filesToRead[1:]
		if _, ok := fileGraph[currentFile]; !ok {

			if _, err := os.Stat(filepath.Join("input", currentFile)); err == nil {
				filesFromCurrentFile, err := readUUIDSFromFile(currentFile)
				if err != nil {
					return nil, err
				}
				fileGraph[currentFile] = filesFromCurrentFile
				filesToRead = append(filesToRead, filesFromCurrentFile...)
			}
		}
	}
	return fileGraph, nil
}

func readUUIDSFromFile(fileUUID string) ([]string, error) {

	path := filepath.Join("input", fileUUID)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}, err
	}

	re := regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
	uuids := re.FindAllString(string(data), -1)

	return uuids, nil
}

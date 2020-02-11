package master

import (
	"crypto/sha256"
	"dister/protos"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type ITestRepository interface {
	AddTask(task *protos.TaskData) (string, error)
	FindTask() ([]*protos.TaskData, error)
	UpdateTask(task *protos.TaskData) error
}

var _ ITestRepository = &MemoryTestRepository{}

type MemoryTestRepository struct {
	v  map[string]*protos.TaskData
	mu sync.Mutex
}

func (rep *MemoryTestRepository) UpdateTask(task *protos.TaskData) error {
	rep.mu.Lock()
	defer rep.mu.Unlock()
	rep.v[task.Id] = task
	return nil
}

func (rep *MemoryTestRepository) AddTask(task *protos.TaskData) (string, error) {
	bytes, err := json.Marshal(task.Task)
	if err != nil {
		return "", err
	}
	id := Sha256Encode(bytes)
	task.Id = id

	rep.mu.Lock()
	defer rep.mu.Unlock()
	rep.v[id] = task
	fmt.Println("add task", task.Id)
	return id, nil
}

func (rep *MemoryTestRepository) FindTask() ([]*protos.TaskData, error) {
	var list []*protos.TaskData
	rep.mu.Lock()
	defer rep.mu.Unlock()

	for _, v := range rep.v {
		if time.Since(time.Unix(v.CreateTime, 0)) >= time.Second*time.Duration(v.Delay) && v.Result == nil {
			list = append(list, v)
		}
	}
	return list, nil
}

func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func Sha256Encode(date []byte) string {
	d := sha256.Sum256(date)
	return Base64Encode(d[:])
}

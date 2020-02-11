package master

import (
	"dister/protos"
	"time"
)

type TaskService struct {
	rep ITaskRepository
}

func (svc *TaskService) AddTask(req *protos.TaskData) (string, error) {
	req.CreateTime = time.Now().Unix()
	return svc.rep.AddTask(req)
}

func (svc *TaskService) FindTask() ([]*protos.TaskData, error) {
	return svc.rep.FindTask()
}

func (svc *TaskService) UpdateTask(task *protos.TaskData) error {
	return svc.rep.UpdateTask(task)
}

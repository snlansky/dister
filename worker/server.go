package worker

import (
	"context"
	"dister/protos"
	"fmt"
)

type Service struct {
}

func (s *Service) Unit(ctx context.Context, task *protos.TaskData) (*protos.TaskData, error) {
	fmt.Println("start unit")
	task.Result = &protos.Result{
		Code: 200,
		Msg:  "ok",
	}
	return task, nil
}

func (s *Service) Prepare(ctx context.Context, req *protos.TaskData) (*protos.TaskProcessResponse, error) {
	panic("implement me")
}

func (s *Service) Commit(ctx context.Context, req *protos.TaskCommitRequest) (*protos.TaskCommitResponse, error) {
	fmt.Println(req.Id)
	return nil, nil
}

func (s *Service) State(ctx context.Context, req *protos.StateRequest) (*protos.StateResponse, error) {
	return &protos.StateResponse{
		St:    protos.StateResponse_Idle,
		Tasks: []string{},
	}, nil
}

func NewService() protos.DisterServer {
	return &Service{}
}

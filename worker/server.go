package worker

import (
	"context"
	"dister/protos"
)

type Service struct {
}

func (s *Service) Prepare(context.Context, *protos.TaskProcessRequest) (*protos.TaskProcessResponse, error) {
	panic("implement me")
}

func (s *Service) Commit(context.Context, *protos.TaskCommitRequest) (*protos.TaskCommitResponse, error) {
	panic("implement me")
}

func (s *Service) State(context.Context, *protos.StateRequest) (*protos.StateResponse, error) {
	panic("implement me")
}

func NewService() protos.DisterServer {
	return &Service{}
}

package master

import (
	"dister/protos"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"log"
)

type Server struct {
	manager *Manager
}

func (s *Server) Register(stream protos.Dister_RegisterServer) error {
	registerMsg, err := stream.Recv()
	if err != nil {
		return err
	}
	if registerMsg.MessageType != protos.Message_REGISTER_REQUEST {
		return errors.New("message type error")
	}
	var r protos.RegisterRequest
	err = proto.Unmarshal(registerMsg.Content, &r)
	if err != nil {
		return errors.WithMessage(err, "proto unmarshal error")
	}

	worker := NewWorker(stream)

	s.manager.AddWorker(worker)
	log.Printf(fmt.Sprintf("worker: %s register\n", worker.id))
	return worker.Start()
}

func NewServer(manager *Manager) protos.DisterServer {
	return &Server{manager: manager}
}

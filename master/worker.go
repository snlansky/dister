package master

import (
	"dister/protos"
	"dister/utils"
	"google.golang.org/grpc"
	"log"
)

type Worker struct {
	id     string
	cpu    int32
	mem    int32
	tasks  []string
	status protos.StateResponse_StatueType
	conn   *grpc.ClientConn
}

func NewWorker(conn *grpc.ClientConn) *Worker {
	return &Worker{
		id:     utils.RandStringRunes(8),
		cpu:    0,
		mem:    0,
		tasks:  []string{},
		status: protos.StateResponse_UnReady,
		conn:   conn,
	}
}

func (w *Worker) Start() error {
	defer func() {
		// todo
		log.Print("closed .....")
	}()
	return nil
}

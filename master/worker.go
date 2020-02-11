package master

import (
	"dister/protos"
	"dister/worker"
	"google.golang.org/grpc"
	"time"
)

type Worker struct {
	id     string
	cpu    int32
	mem    int32
	tasks  []string
	status protos.StateResponse_StatueType
	conn   *grpc.ClientConn
	closeC chan struct{}
}

func NewWorker(id string, conn *grpc.ClientConn) *Worker {
	return &Worker{
		id:     id,
		cpu:    0,
		mem:    0,
		tasks:  []string{},
		status: protos.StateResponse_UnReady,
		conn:   conn,
		closeC: make(chan struct{}),
	}
}

func (w *Worker) Start() error {
	for {
		tick := time.NewTicker(time.Second)
		select {
		case <-tick.C:
			state, err := worker.State(w.conn, &protos.StateRequest{})
			if err != nil {
				return err
			}
			w.tasks = state.Tasks
			w.status = state.St
		case <-w.closeC:
			return nil
		}
	}
}

func (w *Worker) UnitTest(req *protos.TaskData) (*protos.TaskData, error) {
	return worker.UnitTest(w.conn, req)
}

func (w *Worker) Prepare(req *protos.TaskData) (*protos.TaskProcessResponse, error) {
	return worker.Prepare(w.conn, req)
}

func (w *Worker) Commit(req *protos.TaskCommitRequest) (*protos.TaskCommitResponse, error) {
	return worker.Commit(w.conn, req)
}

func (w *Worker) Close() {
	close(w.closeC)
	w.conn.Close()
}

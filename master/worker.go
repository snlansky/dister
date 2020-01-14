package master

import (
	"dister/handler"
	"dister/protos"
	"dister/utils"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"log"
)

var _ handler.MessageHandler = &Worker{}

type Worker struct {
	id      string
	cpu     int32
	mem     int32
	tasks   []string
	isLogin *atomic.Bool
	status  protos.PingRequest_StatueType
	stream  protos.Dister_RegisterServer
}

func (w *Worker) RegisterRequest(msg *protos.Message, request *protos.RegisterRequest) error {
	if w.isLogin.Load() {
		return nil
	}

	w.cpu = request.Cpu
	w.mem = request.Mem
	w.isLogin.Store(true)

	resp := protos.RegisterResponse{
		Id: w.id,
	}
	bytes, err := proto.Marshal(&resp)
	if err != nil {
		return errors.WithMessage(err, "proto marshal error")
	}

	err = w.stream.Send(&protos.Message{
		MessageType:   protos.Message_REGISTER_RESPONSE,
		CorrelationId: msg.CorrelationId,
		Content:       bytes,
	})
	if err != nil {
		return errors.WithMessage(err, "send message error")
	}

	log.Printf("worker: %s, cpu: %d, mem: %d connected\n", w.id, w.cpu, w.mem)
	return nil
}

func (w *Worker) RegisterResponse(msg *protos.Message, response *protos.RegisterResponse) error {
	return errors.New("unexpect type")
}

func (w *Worker) TaskProcessRequest(msg *protos.Message, request *protos.TaskProcessRequest) error {
	return errors.New("unexpect type")
}

func (w *Worker) TaskProcessResponse(msg *protos.Message, response *protos.TaskProcessResponse) error {
	if !w.isLogin.Load() {
		return errors.New("unlogin")
	}
	panic("implement me")
}

func (w *Worker) PingRequest(msg *protos.Message, request *protos.PingRequest) error {
	if !w.isLogin.Load() {
		return errors.New("unlogin")
	}
	w.status = request.St
	w.tasks = request.Tasks
	return nil
}

func (w *Worker) PingResponse() error {
	return errors.New("unexpect type")
}

func NewWorker(stream protos.Dister_RegisterServer) *Worker {
	return &Worker{
		id:      utils.RandStringRunes(8),
		cpu:     0,
		mem:     0,
		tasks:   []string{},
		stream:  stream,
		status:  protos.PingRequest_UnReady,
		isLogin: atomic.NewBool(false),
	}
}

func (w *Worker) Start() error {
	defer func() {
		// todo
		log.Print("closed .....")
	}()
	for {
		msg, err := w.stream.Recv()
		if err != nil {
			return errors.WithMessage(err, "recv message error")
		}
		err = handler.Handler(msg, w)
		if err != nil {
			return errors.WithMessage(err, "handler message error")
		}
	}
}

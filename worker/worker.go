package worker

import (
	"dister/handler"
	"dister/master"
	"dister/protos"
	"dister/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"log"
)

var _ handler.MessageHandler = &Worker{}

func Start(c *cli.Context) error {
	consul := c.String("consul")
	fmt.Println(consul)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(c.String("master_address"), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client, err := master.Register(conn)
	if err != nil {
		return err
	}
	w := NewWorker(client)
	return w.Start()
}

type Worker struct {
	client protos.Dister_RegisterClient
	id     *atomic.String
}

func NewWorker(client protos.Dister_RegisterClient) *Worker {
	return &Worker{client: client, id: atomic.NewString("")}
}

func (w *Worker) RegisterRequest(msg *protos.Message, request *protos.RegisterRequest) error {
	panic("unexpect type")
}

func (w *Worker) RegisterResponse(msg *protos.Message, response *protos.RegisterResponse) error {
	w.id.Store(response.Id)
	return nil
}

func (w *Worker) TaskProcessRequest(msg *protos.Message, request *protos.TaskProcessRequest) error {
	panic("implement me")
}

func (w *Worker) TaskProcessResponse(msg *protos.Message, response *protos.TaskProcessResponse) error {
	panic("unexpect type")
}

func (w *Worker) PingRequest(msg *protos.Message, request *protos.PingRequest) error {
	panic("unexpect type")
}

func (w *Worker) PingResponse() error {
	return nil
}

func (w *Worker) Start() error {
	reg := protos.RegisterRequest{
		Cpu: 4,
		Mem: 8,
	}
	bytes, err := proto.Marshal(&reg)
	if err != nil {
		return err
	}

	corid := utils.GenerateId()
	err = w.client.Send(&protos.Message{
		MessageType:   protos.Message_REGISTER_REQUEST,
		CorrelationId: corid,
		Content:       bytes,
	})
	if err != nil {
		return errors.WithMessage(err, "send register information error")
	}

	for {
		msg, err := w.client.Recv()
		if err != nil {
			return err
		}

		err = handler.Handler(msg, w)
		if err != nil {
			return err
		}
	}
}

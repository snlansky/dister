package handler

import (
	"dister/protos"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type MessageHandler interface {
	RegisterRequest(msg *protos.Message, request *protos.RegisterRequest) error
	RegisterResponse(msg *protos.Message, response *protos.RegisterResponse) error
	TaskProcessRequest(msg *protos.Message, request *protos.TaskProcessRequest) error
	TaskProcessResponse(msg *protos.Message, response *protos.TaskProcessResponse) error
	PingRequest(msg *protos.Message, request *protos.PingRequest) error
	PingResponse() error
}

func Handler(msg *protos.Message, handler MessageHandler) error {
	switch msg.MessageType {
	case protos.Message_REGISTER_REQUEST:
		var m protos.RegisterRequest
		err := proto.Unmarshal(msg.Content, &m)
		if err != nil {
			return err
		}
		return handler.RegisterRequest(msg, &m)

	case protos.Message_REGISTER_RESPONSE:
	case protos.Message_TASK_PROCESS_REQUEST:
	case protos.Message_TASK_PROCESS_RESPONSE:
	case protos.Message_PING_REQUEST:
	case protos.Message_PING_RESPONSE:
	default:
		return errors.New("unexpect message type:" + protos.Message_MessageType_name[int32(msg.MessageType)])
	}
	return nil
}

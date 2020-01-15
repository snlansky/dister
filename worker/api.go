package worker

import (
	"context"
	"dister/protos"
	"google.golang.org/grpc"
)

func Prepare(conn *grpc.ClientConn, request *protos.TaskProcessRequest) (*protos.TaskProcessResponse, error) {
	client := protos.NewDisterClient(conn)
	return client.Prepare(context.Background(), request)
}
func Commit(conn *grpc.ClientConn, request *protos.TaskCommitRequest) (*protos.TaskCommitResponse, error) {
	client := protos.NewDisterClient(conn)
	return client.Commit(context.Background(), request)
}

func State(conn *grpc.ClientConn, request *protos.StateRequest) (*protos.StateResponse, error) {
	client := protos.NewDisterClient(conn)
	return client.State(context.Background(), request)
}



package master

import (
	"context"
	"dister/protos"
	"google.golang.org/grpc"
)

func Register(conn *grpc.ClientConn) (protos.Dister_RegisterClient, error) {
	client := protos.NewDisterClient(conn)
	return client.Register(context.Background())
}

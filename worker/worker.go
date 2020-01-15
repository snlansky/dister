package worker

import (
	"dister/pkg/grpcsr"
	"dister/protos"
	"fmt"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

func Start(c *cli.Context) error {
	fatal := make(chan error)

	port := c.Int("grpc_port")

	go func() {
		fatal <- startGRPC(fmt.Sprintf(":%d", port))
	}()

	consul := c.String("consul")
	register, err := grpcsr.NewConsulRegister(consul, "worker", port)
	if err != nil {
		return err
	}

	err = register.Register()
	if err != nil {
		return err
	}
	defer register.Deregister()

	return <-fatal
}

func startGRPC(address string) error {
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("start grpc server:", address)

	srv := grpc.NewServer()
	svc := NewService()

	protos.RegisterDisterServer(srv, svc)
	grpc_health_v1.RegisterHealthServer(srv, &HealthImpl{})

	return srv.Serve(lis)
}


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
	"strconv"
)

func Start(c *cli.Context) error {
	fatal := make(chan error)

	port := c.Int("grpc_port")

	go func() {
		fatal <- startGRPC(fmt.Sprintf(":%d", port))
	}()

	consul := c.String("consul")
	meta := map[string]string{
		"cpu": strconv.Itoa(4),
		"mem": strconv.Itoa(16),
	}
	register, err := grpcsr.NewConsulRegister(consul, "worker", port, nil, meta)
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
	grpc_health_v1.RegisterHealthServer(srv, &grpcsr.HealthImpl{})

	return srv.Serve(lis)
}


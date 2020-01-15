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
	consul := c.String("consul")
	fmt.Println(consul)
	fatal := make(chan error)

	port := c.Int("grpc_port")

	go func() {
		fatal <- startGRPC(fmt.Sprintf(":%d", port))
	}()

	go func() {
		consulAddr := c.String("consul")
		fatal <- startRegister(consulAddr, "worker", port)
	}()

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

func startRegister(consul string, app string, port int) error {
	register := grpcsr.NewConsulRegister(consul, app, port)
	return register.Register()
}

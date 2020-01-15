package worker

import (
	"dister/protos"
	"fmt"
	"github.com/gin-gonic/gin"
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

	go func() {
		fatal <- startGRPC(c.String("grpc_address"))
	}()

	go func() {
		fatal <- startHttp(c.String("http_address"))
	}()

	return <- fatal
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

func startHttp(address string) error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r.Run(address)
}

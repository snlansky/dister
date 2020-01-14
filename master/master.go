package master

import (
	"dister/protos"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
)

func Start(c *cli.Context) error {
	db := c.String("db")
	fmt.Println(db)
	fatal := make(chan error)

	go func() {
		fatal <- startGRPC(c.String("grpc_address"))
	}()

	go func() {
		fatal <- startHttp(c.String("http_address"))
	}()

	go func() {
		fatal <- startCron()
	}()

	return <-fatal
}

func startGRPC(address string) error {
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("start grpc server:", address)

	s := grpc.NewServer()

	manager := NewManager()
	server := NewServer(manager)

	protos.RegisterDisterServer(s, server)
	return s.Serve(lis)
}

func startHttp(address string) error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func startCron() error {
	ticker := time.After(time.Second*2)
	for {
		select {
		case <-ticker:
			fmt.Println("run task ...")
			time.Sleep(time.Second*5)
			ticker = time.After(time.Second*2)
		}
	}
	return nil
}

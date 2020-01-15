package master

import (
	"dister/pkg/grpcsr"
	"dister/protos"
	"dister/worker"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/urfave/cli"

	"google.golang.org/grpc"
)

func Start(c *cli.Context) error {
	db := c.String("db")
	fmt.Println(db)
	fatal := make(chan error)

	go func() {
		fatal <- startHttp(c.String("http_address"))
	}()

	go func() {
		fatal <- startCron()
	}()

	go func() {
		fatal <- startDiscover(c.String("consul"))
	}()

	conn, err := grpc.Dial(
		"consul://127.0.0.1:8500/grpc.health.v1.worker?wait=14s",
		grpc.WithInsecure(),
		grpc.WithBalancerName("round_robin"),
	)
	if err != nil {
		return err
	}
	_, err = worker.Commit(conn, &protos.TaskCommitRequest{
		Id: "fsfsff",
	})
	if err != nil {
		return err
	}

	return <-fatal
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

func startCron() error {
	ticker := time.After(time.Second * 2)
	for {
		select {
		case <-ticker:
			fmt.Println("run task ...")
			time.Sleep(time.Second * 5)
			ticker = time.After(time.Second * 2)
		}
	}
	return nil
}

func startDiscover(cousul string) error {
	resolver := grpcsr.NewConsulResolver(cousul, "grpc.health.v1.worker")
	resolve, err := resolver.Resolve("")
	if err != nil {
		return err
	}

	for {
		list, err := resolve.Next()
		if err != nil {
			return err
		}

		for _, item := range list {
			fmt.Println(item)
		}

	}
}

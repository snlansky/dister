package master

import (
	"dister/pkg/grpcsr"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/snlansky/glibs/logging"
	"github.com/urfave/cli"
)

var logger = logging.MustGetLogger("master")

func Start(c *cli.Context) error {
	db := c.String("db")
	fmt.Println(db)
	fatal := make(chan error)

	man := NewManager()

	go func() {
		fatal <- startHttp(c.String("http_address"))
	}()

	go func() {
		fatal <- startCron(man)
	}()

	go func() {
		fatal <- startDiscover(c.String("consul"), man)
	}()

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

func startCron(manager *Manager) error {
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

func startDiscover(consul string, manager *Manager) error {
	resolver := grpcsr.NewConsulResolver(consul, "grpc.health.v1.worker")
	resolve, err := resolver.Resolve("")
	if err != nil {
		return err
	}

	for {
		list, err := resolve.Next()
		if err != nil {
			return err
		}

		for _, svc := range list {
			switch svc.Op {
			case naming.Add:
				conn, err := grpc.Dial(svc.Addr, grpc.WithInsecure(), grpc.WithBlock())
				if err != nil {
					logger.Errorf("connect worker [%s] error", svc.Addr)
					continue
				}
				worker := NewWorker(svc.Addr, conn)
				manager.AddWorker(worker)
				logger.Infof("add worker [%s]", svc.Addr)
				go func() {
					defer func() {
						manager.DeleteWorker(worker.id)
						worker.Close()
						logger.Infof("delete worker [%s]", svc.Addr)
					}()
					err := worker.Start()
					if err != nil {
						logger.Warnf("worker start error: %v", err)
						return
					}
				}()
			case naming.Delete:
				worker := manager.DeleteWorker(svc.Addr)
				if worker != nil {
					worker.Close()
					logger.Infof("delete worker [%s]", svc.Addr)
				}
			}
		}
	}
}

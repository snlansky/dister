package master

import (
	"dister/pkg/grpcsr"
	"dister/protos"
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
	svc := &TaskService{rep: &MemoryTaskRepository{v: map[string]*protos.TaskData{}}}

	go func() {
		fatal <- startHttp(c.String("http_address"), svc)
	}()

	go func() {
		fatal <- startCron(man, svc)
	}()

	go func() {
		fatal <- startDiscover(c.String("consul"), man)
	}()

	return <-fatal
}

func startHttp(address string, svc *TaskService) error {
	r := gin.Default()

	ctrl := &Controller{svc: svc}
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		api.POST("/test", ctrl.addTest)
	}
	return r.Run(address)
}

func startCron(manager *Manager, svc *TaskService) error {
	ticker := time.After(time.Second * 2)
	for {
		select {
		case <-ticker:
			tasks, err := svc.FindTask()
			if err != nil {
				return err
			}
			for _, t := range tasks {
				go func() {
					if t.Threads == 0 {
						worker := manager.GetRandWorker()
						logger.Infof("start task [%s], at worker %s", t.Id, worker.id)
						task, err := worker.UnitTest(t)
						if err != nil {
							logger.Errorf("run task [%s], at worker [%s], error: %v", t.Id, worker.id, err)
							return
						}
						if task != nil {
							err := svc.UpdateTask(task)
							if err != nil {
								logger.Errorf("update task [%v], error: %v", task, err)
							}
						}
					}
				}()
			}
			ticker = time.After(time.Second * 2)
		}
	}
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

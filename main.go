package main

import (
	"dister/master"
	"dister/worker"
	"log"
	"os"

	"github.com/urfave/cli"
)

//go:generate protoc -I protos/ protos/dister.proto --go_out=plugins=grpc:protos

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:     "master",
			Aliases:  []string{"m"},
			Usage:    "run master server",
			Category: "master",
			Action:   master.Start,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:   "db",
					Value:  "127.0.0.1:3306",
					Usage:  "mysql dsn",
					EnvVar: "DB",
				},
				&cli.StringFlag{
					Name:   "consul",
					Value:  "127.0.0.1:8500",
					Usage:  "register to consul",
					EnvVar: "CONSUL_REGISTER",
				},
				&cli.StringFlag{
					Name:   "http_address",
					Value:  ":8080",
					Usage:  "http address",
					EnvVar: "HTTP_ADDRESS",
				},
			},
		},
		{
			Name:     "worker",
			Aliases:  []string{"w"},
			Usage:    "run worker to test",
			Category: "worker",
			Action:   worker.Start,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:   "consul",
					Value:  "127.0.0.1:8500",
					Usage:  "register to consul",
					EnvVar: "CONSUL_REGISTER",
				},
				&cli.IntFlag{
					Name:   "grpc_port",
					Value:  9090,
					Usage:  "grpc port",
					EnvVar: "GRPC_PORT",
				},
			},
		},
	}

	app.Name = "dister"
	app.Usage = "application usage"
	app.Description = "distribution tester" // 描述
	app.Version = "1.0.0"                   // 版本

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

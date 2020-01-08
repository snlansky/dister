package worker

import (
	"fmt"
	"github.com/urfave/cli"
)

func Start(c *cli.Context) error {
	consul := c.String("consul")
	fmt.Println(consul)
	return nil

}

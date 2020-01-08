package master

import (
	"fmt"
	"github.com/urfave/cli"
)

func Start(c *cli.Context) error {
	db := c.String("db")
	fmt.Println(db)
	return nil
}

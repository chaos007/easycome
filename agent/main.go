package main

import (
	"os"

	"github.com/chaos007/easycome/libs/utils"

	_ "github.com/chaos007/easycome/data/pb"

	_ "github.com/go-sql-driver/mysql"

	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	// to catch all uncaught panic
	defer utils.PrintPanicStack()
	// go http.ListenAndServe("0.0.0.0:6060", nil)

	app := &cli.App{
		Name:    "agent",
		Usage:   "a gateway for games with stream multiplexing",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Value: "agent1",
				Usage: "id of this service",
			},
			&cli.StringFlag{
				Name:  "server-info",
				Value: "/info",
				Usage: "info of all server info",
			},
			&cli.StringSliceFlag{
				Name:  "etcd-hosts",
				Value: cli.NewStringSlice("http://127.0.0.1:2379"),
				Usage: "etcd hosts",
			},
		},
		Action: func(c *cli.Context) error {
			
			startup(c)
			// wait forever
			select {}
		},
	}
	app.Run(os.Args)

}

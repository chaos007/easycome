package main

import (
	// "game/client_handler"

	// "game/kafka"
	// "game/numbers"

	"os"

	_ "github.com/go-sql-driver/mysql"

	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	// go http.ListenAndServe("0.0.0.0:6060", nil)
	app := &cli.App{
		Name:    "game",
		Usage:   "a stream processor based game",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "group",
				Value: "s1",
				Usage: "id of this group",
			},
			&cli.StringFlag{
				Name:  "id",
				Value: "game1",
				Usage: "id of this service",
			},
			&cli.StringFlag{
				Name:  "server-type",
				Value: "game",
				Usage: "type of this service",
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

			select {}
		},
	}
	app.Run(os.Args)
}

package main

import (
	cli "gopkg.in/urfave/cli.v2"
)

func startup(c *cli.Context) {
	go sigHandler()
	go serviceInit(c)
}

package main

import (
	"github.com/chaos007/easycome/libs/session/mixkcp"
	"github.com/chaos007/easycome/libs/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var (
	wg = new(sync.WaitGroup)
	// server close signal
	die = make(chan struct{})

	checkMe = make(chan int)
)

// handle unix signals
func sigHandler() {
	defer utils.PrintPanicStack()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	for {
		select {
		case msg := <-ch:
			switch msg {
			case syscall.SIGTERM: // å…³é—­agent
				close(die)
				mixkcp.CloseAllSession() //关闭当前所有的连接
				log.Info("sigterm received")
				log.Info("waiting for agents close, please wait...")
				wg.Wait()
				log.Info("agent shutdown.")
				os.Exit(0)
			}
		case <-checkMe:
			close(die)
			mixkcp.CloseAllSession() //关闭当前所有的连接
			log.Info("sigterm received")
			log.Info("waiting for agents close, please wait...")
			wg.Wait()
			log.Info("agent shutdown.")
			os.Exit(0)
		}

	}
}

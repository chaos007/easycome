package main

import (
	"github.com/chaos007/easycome/center/handle"
	"github.com/chaos007/easycome/enum"
	"github.com/chaos007/easycome/libs/config"
	"github.com/chaos007/easycome/libs/etcdservices"
	"github.com/chaos007/easycome/libs/session/grpc"
	"github.com/chaos007/easycome/model/configdata"
	"encoding/json"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	// go http.ListenAndServe("0.0.0.0:6060", nil)
	app := &cli.App{
		Name:    "game",
		Usage:   "a stream processor based game",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Value: "center",
				Usage: "id of this service",
			},
			&cli.StringFlag{
				Name:  "server-info",
				Value: "/info",
				Usage: "info of all server info",
			},
			&cli.StringFlag{
				Name:  "config-path",
				Value: "config",
				Usage: "info of design config info",
			},
		},
		Action: func(c *cli.Context) error {
			config.ServerConfigInit()
			centerConfig := config.GetServerConfig()
			log.Println("id:", c.String("id"))
			log.Println("listen:", centerConfig.Center.GrpcListen)
			log.Println("etcd-hosts:", centerConfig.Center.ETCDHost)
			log.Println("etcd-root:", centerConfig.Center.ETCDRoot)
			log.Println("etcd-watch:", centerConfig.Center.ETCDWatch)
			log.Println("etcd-path:", centerConfig.Center.ETCDServicePath)

			level, err := log.ParseLevel(centerConfig.Center.LogLevel)
			if err != nil {
				log.Fatal(err)
				return err
			}
			log.SetLevel(level)

			handle.Init()

			// 生成配置表信息
			_, err = configdata.DesignConfigInit(c.String("config-path"))
			if err != nil {
				log.Fatal(err)
				return err
			}

			// 向etcd注册自己
			go func() {
				serviceInfo := etcdservices.ServiceInfo{IP: centerConfig.Center.GrpcListen}
				bs, err := etcdservices.BeforeService(centerConfig.Center.ETCDHost)
				if err != nil {
					log.Fatal(err)
				}
				bs.SetService(c.String("id"), serviceInfo, centerConfig.Center.ETCDServicePath)

				info, err := json.Marshal(centerConfig)
				if err != nil {
					log.Fatal(err)
					return
				}
				bs.PutKey(c.String("server-info"), string(info), 30)

				// 监控所有服务器
				etcdservices.Init(centerConfig.Center.ETCDHost, centerConfig.Center.ETCDRoot,
					centerConfig.Center.ETCDWatch,
					bs)

				err = bs.Start()
				log.Fatalln("start keep live error:", err)
				os.Exit(-1)
			}()

			conf := &grpc.Config{
				ServerType:   enum.ServerTypeCenter,
				ServerKey:    centerConfig.Center.ETCDServicePath + c.String("id"),
				Listen:       centerConfig.Center.GrpcListen,
				ReadDeadline: time.Duration(centerConfig.General.ReadDeadLine) * time.Second,
				Sockbuf:      centerConfig.General.SockBuf,
				UDPSockbuf:   centerConfig.General.UDPSockBuf,
				Txqueuelen:   centerConfig.General.QueueLen,
				Dscp:         centerConfig.General.Dscp,
				Sndwnd:       centerConfig.General.UDPSndWnd,
				Rcvwnd:       centerConfig.General.UDPRcvWnd,
				Mtu:          centerConfig.General.UDPMtu,
				Nodelay:      centerConfig.General.Nodelay,
				Interval:     centerConfig.General.Interval,
				Resend:       centerConfig.General.ReSend,
				Nc:           centerConfig.General.Nc,
				WaitGroup:    new(sync.WaitGroup),
			}

			grpc.SetConfig(conf)

			//开启rpc服务，接受内部服务器信息
			go grpc.StartServer()

			go HTTPInit()

			// wait forever
			select {}
		},
	}
	app.Run(os.Args)
}

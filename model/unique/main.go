package main

import (
	"sync"
	"time"

	"github.com/chaos007/easycome/libs/config"
	"github.com/chaos007/easycome/libs/etcdservices"
	"github.com/chaos007/easycome/libs/session/grpc"
	"github.com/chaos007/easycome/unique/handler"

	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:    "game",
		Usage:   "a stream processor based game",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Value: "unique",
				Usage: "id of this service",
			},
			&cli.StringSliceFlag{
				Name:  "etcd-hosts",
				Value: cli.NewStringSlice("http://127.0.0.1:2379"),
				Usage: "etcd hosts",
			},
			&cli.StringFlag{
				Name:  "etcd-root",
				Value: "/backends",
				Usage: "etcd root path",
			},
			&cli.StringFlag{
				Name:  "etcd-service-path",
				Value: "/backends/login/",
				Usage: "etcd service path",
			},
			&cli.StringFlag{
				Name:  "server-type",
				Value: "login",
				Usage: "type of this service",
			},
			&cli.StringSliceFlag{
				Name:  "services",
				Value: cli.NewStringSlice("snowflake"),
				Usage: "auto-discovering services",
			},
		},
		Action: func(c *cli.Context) error {
			config.ServerConfigInit()
			centerConfig := config.GetServerConfig()

			log.Println("id:", c.String("id"))
			log.Println("listen:", c.String("listen"))
			log.Println("etcd-hosts:", c.StringSlice("etcd-hosts"))
			log.Println("etcd-root:", c.String("etcd-root"))
			log.Println("services:", c.StringSlice("services"))
			log.Println("loginPort:", centerConfig.Unique.WebListen)

			level, err := log.ParseLevel(centerConfig.Unique.LogLevel)
			if err != nil {
				log.Fatalln("parse log level error:", err)
				os.Exit(-1)
				return nil
			}
			log.SetLevel(level)

			// 向etcd注册自己
			go func() {
				serviceInfo := etcdservices.ServiceInfo{IP: centerConfig.Unique.GrpcListen}
				bs, err := etcdservices.BeforeService(centerConfig.Unique.ETCDHost)
				if err != nil {
					log.Fatal(err)
				}
				bs.SetService(c.String("id"), serviceInfo, centerConfig.Unique.ETCDServicePath)

				// 监控所有服务器
				etcdservices.Init(centerConfig.Unique.ETCDHost, centerConfig.Unique.ETCDRoot,
					centerConfig.Unique.ETCDWatch, bs)

				err = bs.Start()
				log.Fatalln("start keep live error:", err)
				os.Exit(-1)
			}()

			conf := &grpc.Config{
				ServerType:   c.String("server-type"),
				ServerKey:    centerConfig.Unique.ETCDServicePath + c.String("id"),
				Listen:       centerConfig.Unique.GrpcListen,
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

			handler.Init()
			handler.ProtocolInit()
			//连接数据库

			tabelInit()

			http.HandleFunc("/login", handler.LoginHandler)
			http.HandleFunc("/register", handler.RegisterHandler)
			http.ListenAndServe(centerConfig.Unique.WebListen, nil)
			return nil
		},
	}
	app.Run(os.Args)
}

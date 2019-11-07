package main

import (
	"net/http"
	"os"
	"time"

	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/libs/enum"
	"github.com/chaos007/easycome/game/handle"
	"github.com/chaos007/easycome/libs/etcdservices"
	"github.com/chaos007/easycome/libs/session/grpc"

	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

func serviceInit(c *cli.Context) {
	// 连接etcd
	bs, err := etcdservices.BeforeService(c.StringSlice("etcd-hosts"))
	if err != nil {
		log.Fatalln("before connect to etcd error")
		os.Exit(-1)
		return
	}

	//验证自己是否在服务器列表中
	sType, isMe := bs.CheckMe(c.String("server-info"), c.String("id"), enum.ServerTypeGame)
	if !isMe {
		log.Fatalln("is not me or addr not enough")
		os.Exit(-1)
		return
	}

	level, err := log.ParseLevel(sType.LogLevel)
	if err != nil {
		log.Fatalln("parse log level error:", err)
		os.Exit(-1)
		return
	}

	go func() {
		log.Println(http.ListenAndServe(sType.PProfListen, nil))
	}()

	// 自己发生改变，将自己退出
	go func() {
		bs.WatchMeChange(checkMe)
	}()

	// 向etcd注册自己
	go func() {
		serviceInfo := etcdservices.ServiceInfo{IP: sType.GrpcListen}
		bs.SetService(c.String("id"), serviceInfo, sType.ETCDServicePath)

		// 监控所有服务器
		etcdservices.Init(c.StringSlice("etcd-hosts"), sType.ETCDRoot, sType.ETCDWatch, bs)

		err = bs.Start()
		log.Fatalln("start keep live error:", err)
		os.Exit(-1)
	}()

	// 注册协议
	handle.Init()

	conf := &grpc.Config{
		ServerType:   c.String("server-type"),
		ServerKey:    sType.ETCDServicePath + c.String("id"),
		Listen:       sType.GrpcListen,
		ReadDeadline: time.Duration(bs.General.ReadDeadLine) * time.Second,
		Sockbuf:      bs.General.SockBuf,
		UDPSockbuf:   bs.General.UDPSockBuf,
		Txqueuelen:   bs.General.QueueLen,
		Dscp:         bs.General.Dscp,
		Sndwnd:       bs.General.UDPSndWnd,
		Rcvwnd:       bs.General.UDPRcvWnd,
		Mtu:          bs.General.UDPMtu,
		Nodelay:      bs.General.Nodelay,
		Interval:     bs.General.Interval,
		Resend:       bs.General.ReSend,
		Nc:           bs.General.Nc,
		WaitGroup:    wg,
	}

	grpc.SetConfig(conf)

	// 数据库连接
	// tabelInit(sType.MysqlConfig)

	//开启rpc服务，接受内部服务器信息
	go grpc.StartServer()

	log.Infoln("grpc listen on:", sType.GrpcListen)

	sess := grpc.NewSession()
	err = sess.SendToSteamServer(enum.ServerTypeCenter, &pb.UpGetConfigToCenter{})
	if err != nil {
		log.Errorln("SendToStreamWithServerKey:", err)
	}

	log.SetLevel(level)
}

package main

import (
	"net/http"
	"os"
	"time"

	"github.com/chaos007/easycome/agent/handle"
	"github.com/chaos007/easycome/enum"
	"github.com/chaos007/easycome/libs/etcdservices"
	"github.com/chaos007/easycome/libs/session/grpc"
	"github.com/chaos007/easycome/libs/session/mixkcp"

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
	sType, isMe := bs.CheckMe(c.String("server-info"), c.String("id"), enum.ServerTypeAgent)
	if !isMe {
		log.Fatalln("is not me or addr not enough")
		os.Exit(-1)
		return
	}

	// log level
	level, err := log.ParseLevel(sType.LogLevel)
	if err != nil {
		log.Fatalln("parse log level error:", err)
		os.Exit(-1)
		return
	}

	// pproflisten
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

	// pro
	handle.Init()

	conf := &grpc.Config{
		ServerType:   enum.ServerTypeAgent,
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

	//开启rpc服务，接受内部服务器信息
	go grpc.StartServer()

	log.Println("grpc listen on:", sType.GrpcListen)

	// 开启面向外部的链接
	mconfig := &mixkcp.Config{
		ServerType:   enum.ServerTypeAgent,
		ServerKey:    sType.ETCDServicePath + c.String("id"),
		Listen:       sType.MixListen,
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
		Die:          die,
	}

	//开启tcp服务，接受玩家消息
	go mixkcp.WebSocketServer(mconfig)

	log.Infoln("mixkcp listen on:", sType.MixListen)

	log.SetLevel(level)

}

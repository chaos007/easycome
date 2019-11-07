package grpc

import (
	"errors"
	"net"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	// "game/client_handler"

	"github.com/chaos007/easycome/pb"
)

// 错误
var (
	ErrorIncorrectFrameType = errors.New("incorrect frame type")
	ErrorServiceNotBind     = errors.New("service not bind")
)

// Config 配置
type Config struct {
	ServerType                    string
	ServerKey                     string
	Listen                        string
	EtcdWatch                     string
	ReadDeadline                  time.Duration
	Sockbuf                       int
	UDPSockbuf                    int
	Txqueuelen                    int
	Dscp                          int
	Sndwnd                        int
	Rcvwnd                        int
	Mtu                           int
	Nodelay, Interval, Resend, Nc int
	WaitGroup                     *sync.WaitGroup
}

var serverConfig *Config

type server struct{}

// SetConfig SetConfig
func SetConfig(c *Config) {
	serverConfig = c
}

// StartServer 开启gprc 服务
func StartServer() {
	if serverConfig == nil {
		log.Errorln("serverConfig do not init")
		os.Exit(-1)
		return
	}
	lis, err := net.Listen("tcp", serverConfig.Listen)
	if err != nil {
		log.Errorln("listen tcp error", err)
		os.Exit(-1)
		return
	}
	log.Info("listening on ", lis.Addr())
	s := grpc.NewServer()
	ins := new(server)
	pb.RegisterServiceServer(s, ins)
	s.Serve(lis)
}

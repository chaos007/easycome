package mixkcp

import (
	"github.com/chaos007/easycome/libs/utils"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	kcp "github.com/xtaci/kcp-go"
	"golang.org/x/net/websocket"
)

// Config 配置
type Config struct {
	ServerType                    string
	ServerKey                     string
	Listen                        string
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
	Die                           chan struct{}
}

var uniqueConfig *Config

// TCPServer tcp链接
func TCPServer(config *Config) {
	if uniqueConfig != nil {
		logrus.Errorln("can not open two mixkcp service")
		os.Exit(-1)
		return
	}
	uniqueConfig = config
	// resolve address & start listening
	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.Listen)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	logrus.Info("tcp server listening on:", listener.Addr())

	// loop accepting
	for {
		select {
		case <-uniqueConfig.Die:
			return
		default:
			conn, err := listener.AcceptTCP()
			if err != nil {
				logrus.Warning("accept failed:", err)
				continue
			}
			// set socket read buffer
			conn.SetReadBuffer(config.Sockbuf)
			// set socket write buffer
			conn.SetWriteBuffer(config.Sockbuf)

			// start a goroutine for every incoming connection for reading
			go handleClient(conn, config)
		}
	}
}

// WebSocketServer websocket
func WebSocketServer(config *Config) {
	if uniqueConfig != nil {
		logrus.Errorln("can not open two mixkcp service")
		os.Exit(-1)
		return
	}
	uniqueConfig = config

	http.Handle("/websocket", websocket.Handler(Echo))
	if err := http.ListenAndServe(config.Listen, nil); err != nil {
		logrus.Errorln("ListenAndServe:", err)
		os.Exit(-1)
		return
	}

	logrus.Info("websocket server listening on:", config.Listen)

}

// Echo Echo
func Echo(ws *websocket.Conn) {

	// loop accepting
	select {
	case <-uniqueConfig.Die:
		return
	default:
		// set socket read buffer
		// ws.SetReadBuffer(uniqueConfig.Sockbuf)
		// set socket write buffer
		// ws.SetWriteBuffer(uniqueConfig.Sockbuf)

		// start a goroutine for every incoming connection for reading
		ws.PayloadType = 0x2
		handleClient(ws, uniqueConfig)
	}
}

// UDPServer kcp连接
func UDPServer(config *Config) {
	if uniqueConfig != nil {
		logrus.Errorln("can not open two mixkcp service")
		os.Exit(-1)
		return
	}
	uniqueConfig = config
	l, err := kcp.Listen(config.Listen)
	checkError(err)
	logrus.Info("udp listening on:", l.Addr())
	lis := l.(*kcp.Listener)

	if err := lis.SetReadBuffer(config.Sockbuf); err != nil {
		logrus.Println("SetReadBuffer", err)
	}
	if err := lis.SetWriteBuffer(config.Sockbuf); err != nil {
		logrus.Println("SetWriteBuffer", err)
	}
	logrus.Info("conifg:", config, config.Dscp)
	if err := lis.SetDSCP(config.Dscp); err != nil {
		logrus.Println("SetDSCP", err)
	}

	// loop accepting
	for {
		select {
		case <-uniqueConfig.Die:
			return
		default:
			conn, err := lis.AcceptKCP()
			if err != nil {
				logrus.Warning("accept failed:", err)
				continue
			}
			// set kcp parameters
			conn.SetWindowSize(config.Sndwnd, config.Rcvwnd)
			conn.SetNoDelay(config.Nodelay, config.Interval, config.Resend, config.Nc)
			conn.SetStreamMode(true)
			conn.SetMtu(config.Mtu)

			// start a goroutine for every incoming connection for reading
			go handleClient(conn, config)
		}
	}
}

func checkError(err error) {
	if err != nil {
		logrus.Fatal(err)
		os.Exit(-1)
	}
}

// PIPELINE #1: handleClient
// the goroutine is used for reading incoming PACKETS
// each packet is defined as :
// | 2B size |     DATA       |
//
func handleClient(conn net.Conn, config *Config) {
	defer utils.PrintPanicStack()
	// 读取2Byte的头信息
	// header := make([]byte, 2)

	// todo session cache
	// 创建一个新的session给这个新的连接
	sess, err := NewAgentSession(conn, config)
	if err != nil {
		logrus.Error("cannot get new session:", err)
		return
	}

	// 入口buff读取
	go sess.InBufferStart(config.ReadDeadline)

	// 出口buff开始
	go sess.OutBufferStart()

	// 开启agent处理进程
	sess.Handle()

}

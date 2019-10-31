package socket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/chaos007/easycome/libs/msgmeta"
	"github.com/chaos007/easycome/tools/testscript/handle"
	"github.com/chaos007/easycome/tools/testscript/player"
	"github.com/chaos007/easycome/tools/testscript/types"

	"github.com/golang/protobuf/proto"
	kcp "github.com/xtaci/kcp-go"
)

// ErrorReadCount ErrorReadCount
var (
	ErrorReadCount = errors.New("read count error")
)

// type
var (
	ConnectServerKCP = "kcp"
	ConnectServerTCP = "tcp"
)

// Socket Socket
type Socket struct {
	host             string
	conn             net.Conn //连接
	readCnt          uint32
	writeCnt         uint8
	autoReconnectSec int64
	readQue          chan []byte
	gw               sync.WaitGroup
	connect          string
}

// SetReconnect SetReconnect
func (s *Socket) SetReconnect(retime int64) {
	s.autoReconnectSec = retime
}

// SetServerType SetServerType
func (s *Socket) SetServerType(t string) {
	s.connect = t
}

// Send Send
func (s *Socket) Send(data proto.Message) (err error) {
	var rawdata []byte
	rawdata, err = msgmeta.BuildPacket(data)
	if err != nil {
		fmt.Printf("send error:%s", err.Error())
		return
	}
	sz := uint16(len(rawdata))
	outbuff := bytes.NewBuffer([]byte{})
	// 写包大小
	if err = binary.Write(outbuff, binary.BigEndian, sz-4); err != nil {
		return
	}
	s.writeCnt++
	// 写序号
	if err = binary.Write(outbuff, binary.BigEndian, s.writeCnt); err != nil {
		return
	}
	if err = binary.Write(outbuff, binary.BigEndian, rawdata); err != nil {
		return
	}
	_, err = s.conn.Write(outbuff.Bytes())
	if err != nil {
		fmt.Println("client send data err:", err)
	}
	// fmt.Printf("client send data:%#v,size:%d\n", data, n)
	return
}

// DoProto DoProto
func (s *Socket) DoProto(in chan []byte, die chan bool, p *player.Player) {
	var err error
	defer func() {
		close(die)
	}()
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				fmt.Printf("do in failed.")
				return
			}
			// fmt.Printf("Doproto got msg:%#v,msglen:%d\n", msg, len(msg))
			payload := bytes.NewReader(msg)

			var ser uint8
			var msgID int32

			// 读取序号
			if err = binary.Read(payload, binary.BigEndian, &ser); err != nil {
				fmt.Println("read ser err:", err)
				return
			}
			//读取协议号
			if err = binary.Read(payload, binary.BigEndian, &msgID); err != nil {
				fmt.Println("read ser msg err:", err)
				return
			}
			rawdata := make([]byte, len(msg)-5)
			if err = binary.Read(payload, binary.BigEndian, &rawdata); err != nil {
				fmt.Println("read rawdata msg err:", err)
				return
			}
			// fmt.Printf("Recv,ser:%d,msgID:%d\n", ser, msgID)
			if ret := handle.Handler(msgID, rawdata, p); ret != nil {
				s.Send(ret)
			}
		}
	}
}

// RecvThread 收的线程
func (s *Socket) RecvThread(sess *types.Session, die chan bool, p *player.Player) (err error) {
	header := make([]byte, 2)
	in := make(chan []byte)
	var n int
	defer func() {
		close(in)
	}()
	go s.DoProto(in, die, p)

	for {
		n, err = io.ReadFull(s.conn, header)
		if err != nil {
			fmt.Printf("read header failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}

		size := binary.BigEndian.Uint16(header)

		payload := make([]byte, size+5)
		_, err = io.ReadFull(s.conn, payload)
		if err != nil {
			fmt.Printf("read payload failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}

		// fmt.Printf("RecvThread,read n:%d,size:%d", n, size)
		//消息投递
		select {
		case in <- payload: // payload queued
		case <-sess.Die: //逻辑层断开socket
			fmt.Printf("connection closed by logic, flag:%v ip:%v", sess.Flag, sess.IP)
			return
		}
	}

}

// Dail 拨号
func (s *Socket) Dail() (err error) {
	if s.host == "" {
		err = errors.New("dail failed")
		return
	}

	if s.connect == ConnectServerTCP {
		// 开始连接
		s.conn, err = net.Dial("tcp", s.host)
		// 连不上
		if err != nil {
			fmt.Printf("#connect%s -> %s", s.host, err.Error())
			return
		}
	} else if s.connect == ConnectServerKCP {
		s.conn, err = kcp.Dial(s.host)
		if err != nil {
			fmt.Printf("#connect%s -> %s", s.host, err.Error())
			return
		}
	}
	return
}

// SetHost SetHost
func (s *Socket) SetHost(host string) {
	s.host = host
}

// NewSocket NewSocket
func NewSocket() *Socket {
	return &Socket{
		readQue: make(chan []byte, 10),
		connect: ConnectServerKCP,
	}
}

package mixkcp

import (
	"encoding/binary"
	"github.com/chaos007/easycome/packet"
	"github.com/chaos007/easycome/utils"
	"io"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// Buffer chan管道 #3: buffer
// 控制客户端的数据包发送
type Buffer struct {
	sendCtrl chan struct{} // 发送关闭信号
	conn     net.Conn      // 连接
	cache    []byte        // 一直读的cache
	// cacheIndex   int
	// isSending    bool
	// sending chan struct{}
	// sendingClose chan struct{} //防止多个写出现错误
}

// rawSendToClient 整包发送
func (b *Buffer) rawSendToClient(data []byte) bool {
	// combine output to reduce syscall.write
	sz := len(data)
	count := 0                                        //客户端验证，当前无用
	binary.BigEndian.PutUint16(b.cache, uint16(sz-4)) //协议id去掉4位
	b.cache[2] = byte(uint8(count))
	copy(b.cache[3:], data)
	log.Debugln("mixkcp send packet to client", data)

	// write data
	n, err := b.conn.Write(b.cache[:sz+3])
	if err != nil {
		log.Warningf("Error send reply data, bytes: %v reason: %v", n, err)
		return false
	}

	return true
}

// func (b *Buffer) copyToCache(data []byte) bool {
// 	sz := len(data)
// 	if b.cacheIndex+sz > packet.PacketLimit+2 {
// 		return false
// 	}

// 	count := 1                                                       //客户端验证，当前无用
// 	binary.BigEndian.PutUint16(b.cache[b.cacheIndex:], uint16(sz-4)) //协议id去掉4位
// 	b.cache[b.cacheIndex+2] = byte(uint8(count))
// 	copy(b.cache[b.cacheIndex+3:], data)
// 	b.cacheIndex += 3 + sz

// 	select {
// 	case <-b.sendingClose:
// 		return false
// 	default:
// 		if b.isSending {
// 			return true
// 		}
// 		b.isSending = true
// 		b.sending <- struct{}{}
// 	}
// 	return true
// }

// func (b *Buffer) send() bool {
// 	for {
// 		copy(b.cacheSend, b.cache[:b.cacheIndex])
// 		index := b.cacheIndex
// 		b.cacheIndex = 0
// 		n, err := b.conn.Write(b.cacheSend[:index])
// 		if err != nil {
// 			log.Warningf("Error send reply data, bytes: %v reason: %v", n, err)
// 			return false
// 		}
// 		if b.cacheIndex == 0 {
// 			break
// 		}
// 	}
// 	log.Debugln("sending over", b.cacheIndex)
// 	return true
// }

// OutBufferStart 开启出口buffer循环
func (b *Buffer) OutBufferStart(outPending chan []byte, outDie chan struct{}) {
	defer utils.PrintPanicStack()
	defer func() {
		close(outDie) //多个生产者，当消费者关闭时，通知关闭
		log.Debugln("------outbuffer close")
	}()
	for {
		select {
		case data, ok := <-outPending:
			if !ok {
				return
			}
			if !b.rawSendToClient(data) {
				log.Debugln("send to client error")
				close(b.sendCtrl) //唯一的生产者，直接自身关闭通知session
				return
			}
			// case <-outDie:
			// 	return
			// case <-b.sending:
			// 	log.Debugln("sending")
			// 	if !b.send() {
			// 		log.Debugln("send to client error")
			// 		close(b.sendCtrl) //唯一的生产者，直接自身关闭通知session
			// 		return
			// 	}
			// 	b.isSending = false
		}
	}
}

// InBufferStart 开启入口buffer循环
func (b *Buffer) InBufferStart(pending chan []byte, d time.Duration) {
	defer utils.PrintPanicStack()
	defer func() {
		close(pending) //唯一的pending生产者，由此关闭
		log.Debugln("------inbuffer close")
	}()
	header := make([]byte, 2)

	for {
		//解决死链接的问题
		b.conn.SetReadDeadline(time.Now().Add(d))

		_, err := io.ReadFull(b.conn, header) //conn关闭来回收携程
		if err != nil {
			log.Errorln("read head err :", err)
			return
		}
		size := binary.BigEndian.Uint16(header)
		payload := make([]byte, size+5)
		_, err = io.ReadFull(b.conn, payload)
		if err != nil {
			log.Errorln("read payload error:", err)
			return
		}
		select {
		case pending <- payload:
		default:
			log.Errorln("in chan full") //入口满了掐掉
			return
		}

	}
}

// NewBuffer 给session创建一个关联的写入缓冲区
func NewBuffer(conn net.Conn, ctrl chan struct{}) *Buffer {
	buf := &Buffer{conn: conn}
	buf.sendCtrl = ctrl
	buf.cache = make([]byte, packet.PacketLimit+2)
	// buf.cacheSend = make([]byte, packet.PacketLimit+2)
	// buf.sending = make(chan struct{})
	// buf.sendingClose = make(chan struct{})
	return buf
}

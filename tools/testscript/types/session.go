package types

import (
	"crypto/rc4"
	"net"
	"time"
)

// Session Session
type Session struct {
	IP      net.IP
	Encoder *rc4.Cipher // 加密器
	Decoder *rc4.Cipher // 解密器
	UserID  int64       // 玩家ID
	GSID    string      // 游戏服ID;e.g.: game1,game2

	Die chan struct{} // 会话关闭信号

	Step chan interface{}
	// 会话标记
	Flag int32

	// 时间相关
	ConnectTime    time.Time // TCP链接建立时间
	PacketTime     time.Time // 当前包的到达时间
	LastPacketTime time.Time // 前一个包到达时间

	// RPS控制
	PacketCount uint32 // 对收到的包进行计数，避免恶意发包
}

// NewSession NewSession
func NewSession() *Session {
	return &Session{
		Die:  make(chan struct{}),
		Step: make(chan interface{}),
	}
}

// Substr 截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

package mixkcp

import (
	"crypto/rc4"
	"errors"
	"net"
	"time"

	"github.com/chaos007/easycome/enum"
	"github.com/chaos007/easycome/msgmeta"
	"github.com/chaos007/easycome/packet"
	"github.com/chaos007/easycome/pb"
	"github.com/chaos007/easycome/session/grpc"
	streamclient "github.com/chaos007/easycome/session/rpc_client"
	"github.com/chaos007/easycome/utils"

	// "github.com/chaos007/easycome/model/player"

	"github.com/golang/protobuf/proto"

	log "github.com/sirupsen/logrus"
)

var sessionCountID int64

// Session 会话
type Session struct {
	MQ         chan pb.Frame // 返回给客户端的异步消息
	Encoder    *rc4.Cipher   // 加密器
	Decoder    *rc4.Cipher   // 解密器
	UserID     string        // sessID
	InDie      chan struct{} // 入口Buff会话关闭信号
	outSideDie chan struct{} // 外部通知会话关闭信号
	mqClose    chan struct{} // 防止多个stream读引起错误
	// OutDie     chan struct{} // 防止多个协程写引起错误

	outPending      chan []byte
	outPendingClose chan struct{}
	// outStreamDie    chan struct{} // rpc发送错误的通知

	inPending chan []byte // 客户端到session数据包chan

	// 会话标记
	Flag      int32
	SessionID int64

	// 时间相关
	ConnectTime     time.Time // TCP链接建立时间
	PacketTime      time.Time // 当前包的到达时间
	LastPacketTime  time.Time // 前一个包到达时间
	PacketCount     int8      // 对收到的包进行计数，避免恶意发包
	PacketCount1Min int       // 每分钟的包统计，用于RPM判断

	sendClientBuff *Buffer
	clientInBuff   *Buffer
	conn           net.Conn
	IP             string
	port           string
	serverType     string //当前会话的类型
	serverKey      string //当前的服务器id

	closeCallBackList *methodList

	grpcSession *grpc.Session

	rpcClients *streamclient.RPCStreamMap
}

// NewAgentSession 新agent会话
func NewAgentSession(conn net.Conn, c *Config) (*Session, error) {
	sess := &Session{}
	sess.IP = conn.RemoteAddr().String()
	sess.conn = conn
	sess.InDie = make(chan struct{})
	sess.outSideDie = make(chan struct{})
	sess.mqClose = make(chan struct{})
	// sess.OutDie = make(chan struct{})

	sess.outPending = make(chan []byte, c.Txqueuelen)
	sess.outPendingClose = make(chan struct{})

	sess.closeCallBackList = newMethodList()

	// 初始化 session
	sess.MQ = make(chan pb.Frame, 512)
	//连接时间
	sess.ConnectTime = time.Now()
	//最后收包时间
	sess.LastPacketTime = time.Now()
	sess.inPending = make(chan []byte, c.Txqueuelen)

	// 新建一个 write buffer
	sess.sendClientBuff = NewBuffer(conn, sess.InDie)
	sess.clientInBuff = NewBuffer(conn, sess.InDie)
	sess.rpcClients = streamclient.NewRPCClient()
	sess.serverType = c.ServerType
	sess.serverKey = c.ServerKey

	sessionCountID++
	sess.SessionID = sessionCountID

	return sess, nil
}

// SendToServer TODO 携程优化 发送至游戏服
func (s *Session) SendToServer(serverType string, m proto.Message) error {
	if m == nil || s.Flag&enum.SessKickedOut != 0 {
		return errors.New("nil message or close session")
	}
	ret, err := msgmeta.BuildPacket(m)
	if err != nil {
		log.Errorln("msg proto marshal exced error,msg name:", proto.MessageName(m))
		return nil
	}

	return s.sessionSendToServerSteam(serverType, ret)
}

// SessionLogin 登录到服务器，外部的链接未登录，取不到id
func (s *Session) SessionLogin(userid string) {
	if old := GetUserSession(userid); old != nil && s.SessionID != old.SessionID {
		old.SessionClose()
	}

	s.UserID = userid
	setUserSession(s)
}

func (s *Session) sessionSendToServerSteam(serverType string, data []byte) error {
	stream := s.getRPCStream(serverType)
	if stream == nil {
		log.Debugln("stream null")
		return errors.New("none steam null")
	}
	frame := &pb.Frame{
		Message: data,
	}
	if err := stream.Stream.Send(frame); err != nil {
		log.Errorln("sessionSendToServerSteam stream send error")
		return err
	}
	return nil
}

// getRPCStream 获得到内部服务器的流
func (s *Session) getRPCStream(serverType string) *streamclient.RPCStream {
	path := serverType
	if serverType == enum.ServerTypeAgent || serverType == enum.ServerTypeGame { //需要区租的服务器
		path = serverType
	}
	return s.rpcClients.TryGetRPCStream(path, s.serverKey, s.UserID, s.MQ, s.mqClose)
}

// SendToClientUseProto 外部会话发送至客户端
func (s *Session) SendToClientUseProto(data proto.Message) error {
	// 屏蔽空包
	if data == nil || s.Flag&enum.SessKickedOut != 0 {
		return errors.New("nil message or close session")
	}

	ret, err := msgmeta.BuildPacket(data)
	if err != nil {
		return err
	}

	select {
	case <-s.outPendingClose:
		return nil
	case s.outPending <- ret:
	default:
		log.Debugln("out chan full")
		s.SessionClose()
		return errors.New("nil message or close session")
	}

	// if !s.sendClientBuff.copyToCache(ret) {
	// 	log.Debugln("out chan full")
	// 	s.SessionClose()
	// 	return errors.New("nil message or close session")
	// }
	return nil
}

// sendToClient 外部会话发送至客户端，发生错误直接调用关闭，不用返回错误信息
func (s *Session) sendToClient(data []byte) {
	// 屏蔽空包
	if data == nil || s.Flag&enum.SessKickedOut != 0 {
		return
	}

	select {
	case <-s.outPendingClose:
		return
	case s.outPending <- data:
	default:
		log.Debugln("out chan full")
		s.SessionClose()
		return
	}

	// if !s.sendClientBuff.copyToCache(data) {
	// 	log.Debugln("out chan full")
	// 	s.SessionClose()
	// }
	return

}

// OutBufferStart 开启buffer循环
func (s *Session) OutBufferStart() {
	s.sendClientBuff.OutBufferStart(s.outPending, s.outPendingClose)
}

// InBufferStart 开启buffer循环
func (s *Session) InBufferStart(d time.Duration) {
	s.clientInBuff.InBufferStart(s.inPending, d)
}

// SessionClose 需要外部的关闭，outpending的生产者
func (s *Session) SessionClose() {
	if s.Flag&enum.SessKickedOut != 0 {
		return
	}
	close(s.outSideDie)
}

// RegisterSessionDead 注册session关闭的回调函数
func (s *Session) RegisterSessionDead(callback func(string) error) {
	if s.Flag&enum.SessKickedOut != 0 {
		return
	}
	s.closeCallBackList.addMethod(callback)
}

// Close session关闭
func (s *Session) close() {
	log.Debugf("client close ip:%s,port:%s", s.IP, s.port)
	delSession(s.SessionID)
	close(s.outPending)
	// close(s.OutDie)
	s.conn.Close()
	close(s.mqClose) //关闭所有stream的接受通道
	s.rpcClients.Close()
	s.closeCallBackList.exce(s.UserID)
}

// Handle 处理请求
func (s *Session) Handle() {
	defer s.close()
	uniqueConfig.WaitGroup.Add(1)
	defer uniqueConfig.WaitGroup.Done() // wg 减 1, 用于手动关闭服务
	defer utils.PrintPanicStack()

	// minute timer
	minTimer := time.After(time.Minute)

	log.Debugf("new client by ip:%s", s.IP)
	// >> the main message loop <<
	// 处理4种消息
	//  1. 来自客户端的消息
	//  2. 来自其他服的消息
	//  3. timer
	//  4. 关闭服务信号
	for {
		select {
		case msg, ok := <-s.inPending: // inpending的消费者，关闭整个session关闭
			if !ok {
				log.Debugln("receive msg from client error")
				s.Flag |= enum.SessKickedOut //主协程关闭
			} else {
				log.Debugln("agent() receive packet from client:", msg)
				s.PacketCount++
				s.PacketCount1Min++
				s.PacketTime = time.Now()
				s.LastPacketTime = s.PacketTime
				if result := s.route(msg); result != nil {
					if err := s.sessionSendToServerUseProto(result); err != nil { //内部rpc错误不用关闭
						log.Errorln("sessionSendToServerUseProto error", err)
					}
				}
			}

		case frame, ok := <-s.MQ: // packets from other server
			log.Debugln("-----receive msg from other server msg")
			if !ok {
				log.Errorln("MQ close in no way，whole session close")
				s.Flag |= enum.SessKickedOut //主协程关闭
			} else if result := s.routeRPC(frame.Message); result != nil {
				if err := s.sessionSendToServerUseProto(result); err != nil { //内部rpc错误不用关闭
					log.Errorln("route rpc msg use proto error", err)
				}
			}
		case <-minTimer: // minutes timer
			log.Debugln("agent() receive packet from timer")
		case <-s.InDie: // server is shuting down...
			log.Debugln("receive from outbuffer shuting down") //出口buffer发生错误
			s.Flag |= enum.SessKickedOut
		case <-s.outSideDie: // server is shuting down...
			log.Debugln("receive from OutSide shuting down")
			s.Flag |= enum.SessKickedOut
		}
		if s.Flag&enum.SessKickedOut != 0 { //所有协程被回收
			log.Debugln("handle over")
			return
		}
	}
}

// route client protocol
func (s *Session) route(p []byte) proto.Message {
	// start := time.Now()
	defer utils.PrintPanicStack(s, p)

	// 解密
	// if s.Flag&SessEncrypt != 0 {
	// 	s.Decoder.XORKeyStream(p, p)
	// }

	// 封装为reader
	reader := packet.Reader(p)
	// 读客户端数据包序列号(1,2,3...)
	// 客户端发送的数据包必须包含一个自增的序号，必须严格递增
	seqID, err := reader.ReadS8()
	if err != nil || seqID != s.PacketCount {
		log.Error("read client timestamp failed:", err)
		s.Flag |= enum.SessKickedOut //close all session
		return nil
	}

	// 读协议号
	b, err := reader.ReadS32()
	if err != nil {
		log.Error("read protocol number failed.")
		s.Flag |= enum.SessKickedOut //close all session
		return nil
	}
	// 协议号的划分在消息注册模块, 用户可以自定义注册模块，用于转发到不同的后端服务
	log.Debug("route receive msg id:", b)

	serverType := msgmeta.GetMsgServerType(b)

	if serverType != s.serverType && serverType != "" { //不是本服的消息
		if err := s.sessionSendToServerSteam(serverType, p[1:]); err != nil {
			log.Errorf("service id:%v execute failed, error:%v", b, err)
			return nil
		}
	} else { //本服消息协议，进行解析
		data, err := msgmeta.ParsePacket(b, p[5:])
		if err != nil {
			log.Errorf("msg parse packet error,msgid:%d err:%s", b, err)
			return nil
		}
		if method := getProtocol(b); method != nil {
			retProto, err := method(s, data)
			if err != nil {
				log.Errorf("msg exced error,msgid:%d, err:%s", b, err)
				return nil
			}
			return retProto
		}
		log.Errorln("mixkcp can not find resgister method,msgid:", b)
	}

	return nil
}

// route client protocol
func (s *Session) routeRPC(p []byte) proto.Message {
	// start := time.Now()
	defer utils.PrintPanicStack(s, p)

	// 解密
	// if s.Flag&SessEncrypt != 0 {
	// 	s.Decoder.XORKeyStream(p, p)
	// }

	// 封装为reader
	reader := packet.Reader(p)

	// 读协议号
	b, err := reader.ReadS32()
	if err != nil {
		log.Error("read protocol number failed.")
		return nil
	}
	// 协议号的划分在消息注册模块, 用户可以自定义注册模块，用于转发到不同的后端服务
	log.Debug("route receive msg id:", b)

	serverType := msgmeta.GetMsgServerType(b)

	if serverType == "" { //返回客户端的消息
		s.sendToClient(p)
	} else if serverType != s.serverType && serverType != enum.ServerTypeAll { //不是本服的消息
		if err := s.sessionSendToServerSteam(serverType, p); err != nil {
			log.Errorf("service id:%v execute failed, error:%v", b, err)
			return nil
		}
	} else { //本服消息协议，进行解析
		data, err := msgmeta.ParsePacket(b, p[4:])
		if err != nil {
			log.Errorf("msg parse packet error,msgid:%d err:%s", b, err)
			return nil
		}
		if method := getProtocol(b); method != nil {
			retProto, err := method(s, data)
			if err != nil {
				log.Errorf("msg exced error,msgid:%d, err:%s", b, err)
				return nil
			}
			return retProto
		}
		log.Errorln("rpc can not find resgister method,msgid:", b)
	}

	return nil
}

// sessionSendToServerUseProto 用proto划分协议发送至不同服务器
func (s *Session) sessionSendToServerUseProto(data proto.Message) error {
	if data == nil {
		return nil
	}
	id, ok := msgmeta.MessageMetaByName(proto.MessageName(data))
	if !ok {
		log.Errorln("can not find msgid return,msgid:", id)
		return nil
	}
	serverType := msgmeta.GetMsgServerType(id)
	ret, err := msgmeta.BuildPacket(data)
	if err != nil {
		log.Errorln("msg proto marshal exced error,msg name:", proto.MessageName(data))
	}
	if serverType == "" { //返回客户端的消息
		s.sendToClient(ret)
	} else if serverType != s.serverType && serverType != enum.ServerTypeAll { //不是本服的消息
		if err := s.sessionSendToServerSteam(serverType, ret); err != nil {
			log.Errorf("service id:%v execute failed, error:%v", id, err)
			return err
		}
	} else { //本服消息协议，不进行发送
		return errors.New("send message to me server")
	}

	return nil
}

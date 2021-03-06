package grpc

import (
	"path/filepath"

	"github.com/chaos007/easycome/utils"
	"github.com/chaos007/easycome/msgmeta"
	"github.com/chaos007/easycome/packet"
	"github.com/chaos007/easycome/pb"
	streamclient "github.com/chaos007/easycome/session/rpc_client"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// Session 会话:
// 会话是一个单独玩家的上下文，在连入后到退出前的整个生命周期内存在
// 根据业务自行扩展上下文
type Session struct {
	Flag       int32 // 会话标记
	UserID     string
	serverType string
	serverKey  string

	clientServerType string
	ServerStream     pb.Service_StreamServer

	outSideDie chan struct{}

	inFrame           chan *pb.Frame
	closeCallBackList *methodList
	rpcClients        *streamclient.RPCStreamMap

	MQ      chan pb.Frame // 返回给其他服务器的异步消息
	mqClose chan struct{} // 防止多个stream读引起错误

	Data interface{}

	// Player interfacer.Player
}

// SessionClose 需要外部的关闭
func (s *Session) SessionClose() {
	if s.Flag&SessKickedOut != 0 {
		return
	}
	close(s.outSideDie)
}

// RegisterSessionDead 注册session关闭的回调函数
func (s *Session) RegisterSessionDead(callback func(string) error) {
	if s.Flag&SessKickedOut != 0 {
		return
	}
	s.closeCallBackList.addMethod(callback)
}

// SendToStreamWithServerKey 使用服务器id发送到服务器上
func (s *Session) SendToStreamWithServerKey(key string, m proto.Message) error {
	// 屏蔽空包
	if m == nil || s.Flag&SessKickedOut != 0 {
		log.Errorln("SendToStreamWithServerKey close or empty proto message")
		return nil
	}

	stream := s.getRPCStreamWithID(key)
	if stream == nil {
		log.Errorln("cannot open stream: id:", key)
		return nil
	}
	ret, err := msgmeta.BuildPacket(m)
	if err != nil {
		log.Errorln("msg proto marshal exced error,msgname:", proto.MessageName(m))
		return nil
	}

	log.Debugln("Send To Steam Server: id", key)
	frame := &pb.Frame{
		Message: ret,
	}
	if err := stream.Stream.Send(frame); err != nil {
		log.Errorln("sessionSendToServerSteam stream send error")
		return err
	}
	return nil
}

// SendToClientAsServer 发向当前连接的客户端
func (s *Session) SendToClientAsServer(m proto.Message) error {
	// 屏蔽空包
	if m == nil || s.Flag&SessKickedOut != 0 {
		log.Errorln("SendToClientAsServer close or empty proto message")
		return nil
	}
	ret, err := msgmeta.BuildPacket(m)
	if err != nil {
		log.Errorln("msg proto marshal exced error,msgname:", proto.MessageName(m))
		return nil
	}

	frame := &pb.Frame{
		Message: ret,
	}
	if err := s.ServerStream.Send(frame); err != nil {
		log.Errorln("sessionSendToServerSteam stream send error as server")
		return err
	}
	return nil
}

// SendToSteamServerBalance 作为客户端发送到其他服务器，随机发送
func (s *Session) SendToSteamServerBalance(serverType string, m proto.Message, isBalance bool) error {
	// 屏蔽空包
	if m == nil || s.Flag&SessKickedOut != 0 {
		log.Errorln("SendToSteamServerBalance close or empty proto message")
		return nil
	}
	ret, err := msgmeta.BuildPacket(m)
	if err != nil {
		log.Errorln("msg proto marshal exced error,msgname:", proto.MessageName(m))
		return nil
	}

	log.Debugln("Send To Steam Server:", serverType)
	frame := &pb.Frame{
		Message: ret,
	}
	if serverType == s.clientServerType { //跟本次客户端会话一致
		if err := s.ServerStream.Send(frame); err != nil {
			log.Errorln("SendToSteamServerBalance stream send error as server")
			return err
		}
		return nil
	}

	var stream *streamclient.RPCStream

	if isBalance {
		stream = s.getRandomRPCStream(serverType)
	} else {
		stream = s.getRPCStream(serverType)
	}
	if stream == nil {
		log.Errorln("cannot open stream: serverType:", serverType)
		return nil
	}

	if err := stream.Stream.Send(frame); err != nil {
		log.Errorln("SendToSteamServerBalance stream send error:", serverType, proto.MessageName(m))
		return err
	}
	return nil
}

func (s *Session) getRandomRPCStream(serverType string) *streamclient.RPCStream {
	path := serverType
	return s.rpcClients.GetRandomRPCStream(path, s.serverKey, s.UserID, s.MQ, s.mqClose)
}

func (s *Session) getRPCStream(serverType string) *streamclient.RPCStream {
	path := serverType
	return s.rpcClients.TryGetRPCStream(path, s.serverKey, s.UserID, s.MQ, s.mqClose)
}

func (s *Session) getRPCStreamWithID(id string) *streamclient.RPCStream {
	return s.rpcClients.TryGetRPCStreamWithID(id, s.serverKey, s.UserID, s.MQ, s.mqClose)
}

// NewSession 新的session
func NewSession() *Session {
	sess := &Session{}
	sess.outSideDie = make(chan struct{})

	sess.inFrame = make(chan *pb.Frame, 128)
	sess.closeCallBackList = newMethodList()

	sess.rpcClients = streamclient.NewRPCClient()
	sess.MQ = make(chan pb.Frame, 512)
	sess.mqClose = make(chan struct{})
	sess.serverType = serverConfig.ServerType
	sess.serverKey = serverConfig.ServerKey
	return sess
}

// Stream 管道 #2 stream processing
// the center of game logic
func (s *server) Stream(stream pb.Service_StreamServer) error {
	defer utils.PrintPanicStack()
	serverConfig.WaitGroup.Add(1)
	defer serverConfig.WaitGroup.Done()
	log.Debug("stream start")
	//初始化session
	sess := NewSession()
	defer sess.close()

	// 从context中获取数据 TODO FromIncomingContext 确定使用的函数
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Error("cannot read metadata from context")
		return ErrorIncorrectFrameType
	}

	// // 获取userid
	if len(md["userid"]) == 0 {
		log.Errorln("cannot read key:userid from metadata")
		return ErrorIncorrectFrameType
	}
	// // 获取客户端服务器类型
	if len(md["server_key"]) == 0 {
		log.Errorln("cannot read key:server_key from metadata")
		return ErrorIncorrectFrameType
	}

	// 注册 user
	sess.UserID = md["userid"][0]
	sess.clientServerType = filepath.Base(filepath.Dir(md["server_key"][0]))
	sess.ServerStream = stream
	// err := sess.Player.NewPlayer(sess.UserID)
	// if err != nil {
	// 	log.Errorln("user login err:", err)
	// 	return err
	// }
	if sess.UserID != "" { //用户的连接
		SetUserSession(sess)
	} else { //服务器的连接
		setServerSession(sess)
	}
	log.Debugln("flag:", sess.Flag&SessKickedOut == 0)
	log.Debug("userid: ", sess.UserID, " logged in")

	go sess.recv(stream)

	return sess.handle(stream)
}

// OutSideHandle 外部处理
func (s *Session) OutSideHandle(stream pb.Service_StreamServer) error {
	return s.handle(stream)
}

// handle Handle
func (s *Session) handle(stream pb.Service_StreamServer) error {
	for {
		select {
		case msg, ok := <-s.inFrame: // frames from client as server stream
			if !ok || msg == nil { // EOF
				log.Debugln("frames from client as server stream EOF")
				s.Flag |= SessKickedOut
			} else {
				log.Debugln("message recv inFrame", msg)
				if result := s.route(msg.Message); result != nil {
					ret, err := msgmeta.BuildPacket(result)
					if err != nil {
						log.Errorf("msgmeta.BuildPacket err:%s,msgname:%s", err, proto.MessageName(result))
						return nil
					}
					frame := &pb.Frame{
						Message: ret,
					}
					if err := stream.Send(frame); err != nil {
						log.Errorln("sessionSendToServerSteam stream handle send error:", proto.MessageName(msg))
						return err
					}
				}
			}
		case msg, ok := <-s.MQ: // frames from server as client stream
			if !ok { // EOF
				log.Debugln("frames from server as client stream EOF")
				s.Flag |= SessKickedOut
			} else {
				log.Debugln("message recv MQ", msg)
				if result := s.route(msg.Message); result != nil {
					if err := s.sessionSendToServerUseProto(result); err != nil {
						log.Errorln("sessionSendToServerUseProto ,err:", err)
						s.Flag |= SessKickedOut
					}
				}
			}
		case <-s.outSideDie: // server is shuting down...
			log.Debugln("receive from OutSide shuting down")
			s.Flag |= SessKickedOut
		}
		if s.Flag&SessKickedOut != 0 {
			log.Debugln("handle over")
			return nil
		}
	}
}

// 管道 #1 stream receiver
// this function is to make the stream receiving SELECTABLE
func (s *Session) recv(stream pb.Service_StreamServer) {
	log.Debugln("------insteam recv start")
	defer utils.PrintPanicStack()
	defer func() {
		close(s.inFrame) //唯一的生产者
		log.Debugln("------insteam close as server")
	}()
	for {
		in, err := stream.Recv()
		log.Debugln("grpc stream receive msg ------------")
		if err != nil || in == nil { // client closed
			log.Debugln("grpc stream io eof")
			return
		}
		log.Debugln("grpc stream receive", in.Message)
		select {
		case s.inFrame <- in:
		default: //来不及收
			log.Errorln("in chan full")
			return
		}
	}
}

func (s *Session) close() {
	log.Debugln("all session close")
	close(s.mqClose)
	s.rpcClients.Close()
	DelUserSession(s.UserID)
	s.closeCallBackList.exce(s.UserID)
}

func (s *Session) route(p []byte) proto.Message {
	// start := time.Now()
	defer utils.PrintPanicStack(s, p)

	// 封装为reader
	reader := packet.Reader(p)

	// 读协议号
	b, err := reader.ReadS32()
	if err != nil {
		log.Errorln("read protocol number failed. msgid:", b)
		return nil
	}
	// 协议号的划分在消息注册模块, 用户可以自定义注册模块，用于转发到不同的后端服务
	log.Debug("route receive msg id:", b)

	data, err := msgmeta.ParsePacket(b, p[4:])
	if err != nil {
		log.Errorf("msg parse packet error,msgid:%d err:%s", b, err)
		return nil
	}

	serverType := msgmeta.GetMsgServerType(b)

	if serverType == "" { //返回客户端的消息，rpc的消息不能直接返回客户端
		log.Debugln("rpc recv from client error")
		return nil
	} else if serverType != s.serverType && serverType != "ToAll" { //不是本服的消息
		if err := s.SendToSteamServerBalance(serverType, data, false); err != nil {
			log.Errorf("message id:%v execute failed, error:%v", b, err)
			return nil
		}
	} else { //本服消息协议，进行解析
		if method := getProtocol(b); method != nil {
			retProto, err := method(s, data)
			if err != nil {
				log.Errorln("msg exced error,msgid:", b)
				return nil
			}
			return retProto
		}
		log.Errorln("grpc msg id cannot find protocol error,msgid:", b)
		return nil
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
	if serverType == "" { //返回客户端的消息
		return s.SendToClientAsServer(data)
	} else if serverType != s.serverType && serverType != "ToAll" { //不是本服的消息
		if err := s.SendToSteamServerBalance(serverType, data, false); err != nil {
			log.Errorf("message id:%v execute failed, error:%v", id, err)
			return err
		}
	} else { //本服消息协议，不进行发送
		log.Errorln("rpc message can not send to myself")
		return nil
	}

	return nil
}

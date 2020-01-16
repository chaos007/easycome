package streamclient

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"sync"

	"gitlab.dianchu.cc/2020008/libs/pb"

	"io"

	"gitlab.dianchu.cc/2020008/libs/services"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// RPCStreamMap rpcstream列表
type RPCStreamMap struct {
	list map[string]map[string]*RPCStream
	lock *sync.RWMutex
}

// NewRPCClient 一个rpc客户端
func NewRPCClient() *RPCStreamMap {
	return &RPCStreamMap{
		list: map[string]map[string]*RPCStream{},
		lock: new(sync.RWMutex),
	}
}

// GetRPCStream 获得服务器类型的链接
func (r *RPCStreamMap) GetRPCStream(serverType string) *RPCStream {
	r.lock.RLock()
	defer r.lock.RUnlock()
	v := r.list[serverType]
	if v == nil {
		return nil
	}
	list := []*RPCStream{}
	for _, item := range v {
		list = append(list, item)
	}
	return list[rand.Intn(len(list))]
}

// Close Close
func (r *RPCStreamMap) Close() {
	for _, item := range r.list {
		for _, sitem := range item {
			sitem.close()
		}
	}
}

// GetRandomRPCStream 随机获得一个rpc连接
func (r *RPCStreamMap) GetRandomRPCStream(serverType, myServerKey string, userid string, fr chan pb.Frame, mqClose chan struct{}) *RPCStream {
	r.lock.Lock()
	defer r.lock.Unlock()

	gameConn, key := services.GetService2(serverType)
	if gameConn == nil {
		log.Debugln("none gprc conn")
		return nil
	}

	if v, ok := r.list[serverType]; ok {
		if c, cok := v[key]; cok && !c.IsClose {
			return c
		}
	}

	log.Debugln("game sever key :", key)
	cli := pb.NewServiceClient(gameConn)
	// 开启到rpc的流
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{
		"userid":     fmt.Sprint(userid),
		"server_key": myServerKey,
	}))
	stream, err := cli.Stream(ctx)
	if err != nil {
		log.Debugln("stream open err:", err)
		return nil
	}
	rpcStream := &RPCStream{
		GSID:   key,
		Stream: stream,
	}
	m := r.list[serverType]
	if m == nil {
		m = map[string]*RPCStream{
			key: rpcStream,
		}
		r.list[serverType] = m
	} else {
		m[key] = rpcStream
	}

	go rpcStream.rpcClientRecv(fr, mqClose)
	return rpcStream
}

// TryGetRPCStream 尝试获得一个rpc连接
func (r *RPCStreamMap) TryGetRPCStream(serverType, myServerKey string, userid string, fr chan pb.Frame, mqClose chan struct{}) *RPCStream {
	r.lock.Lock()
	defer r.lock.Unlock()
	if v, ok := r.list[serverType]; ok && len(v) > 0 {
		list := []*RPCStream{}
		for _, item := range v {
			if !item.IsClose {
				list = append(list, item)
			}
		}
		if len(list) > 0 {
			idx := rand.Intn(len(list))
			return list[idx]
		}
	}
	gameConn, key := services.GetService2(serverType)
	if gameConn == nil {
		log.Debugln("none gprc conn")
		return nil
	}
	log.Debugln("game sever key :", key)
	cli := pb.NewServiceClient(gameConn)
	// 开启到rpc的流
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{
		"userid":     fmt.Sprint(userid),
		"server_key": myServerKey,
	}))
	stream, err := cli.Stream(ctx)
	if err != nil {
		log.Debugln("stream open err:", err)
		return nil
	}
	rpcStream := &RPCStream{
		GSID:   key,
		Stream: stream,
	}
	m := r.list[serverType]
	if m == nil {
		m = map[string]*RPCStream{
			key: rpcStream,
		}
		r.list[serverType] = m
	} else {
		m[key] = rpcStream
	}

	go rpcStream.rpcClientRecv(fr, mqClose)
	return rpcStream
}

// TryGetRPCStreamWithID 尝试获得一个rpc连接
func (r *RPCStreamMap) TryGetRPCStreamWithID(key, myServerKey string, userid string, fr chan pb.Frame, mqClose chan struct{}) *RPCStream {
	r.lock.Lock()
	defer r.lock.Unlock()

	path := filepath.Base(filepath.Dir(key))

	if v, ok := r.list[path]; ok {
		if r, rok := v[key]; rok && !r.IsClose {
			return r
		}
	}

	gameConn := services.GetServiceWithID(path, filepath.Base(key))
	if gameConn == nil {
		log.Debugln("none gprc conn")
		return nil
	}
	cli := pb.NewServiceClient(gameConn)
	// 开启到rpc的流
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{
		"userid":     fmt.Sprint(userid),
		"server_key": myServerKey,
	}))
	stream, err := cli.Stream(ctx)
	if err != nil {
		log.Debugln("stream open err:", err)
		return nil
	}
	rpcStream := &RPCStream{
		GSID:   key,
		Stream: stream,
	}
	m := r.list[path]
	if m == nil {
		m = map[string]*RPCStream{
			key: rpcStream,
		}
		r.list[path] = m
	} else {
		m[key] = rpcStream
	}

	go rpcStream.rpcClientRecv(fr, mqClose)
	return rpcStream
}

// RPCStreamInit rpcstream初始化
func (r *RPCStreamMap) RPCStreamInit(serverType string, userid int64, fr chan pb.Frame, mqClose chan struct{}) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	gameConn, key := services.GetService2(serverType)
	if gameConn == nil {
		log.Debugln("none gprc conn:", serverType)
		return errors.New("none gprc conn")
	}
	log.Debugln("game sever key :", key)
	cli := pb.NewServiceClient(gameConn)
	// 开启到rpc的流
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{"userid": fmt.Sprint(userid)}))
	stream, err := cli.Stream(ctx)
	if err != nil {
		log.Debugln("stream open err:", err)
		return err
	}
	rpcStream := &RPCStream{
		GSID:   key,
		Stream: stream,
	}
	m := r.list[serverType]
	if m == nil {
		m = map[string]*RPCStream{
			key: rpcStream,
		}
		r.list[serverType] = m
	} else {
		m[key] = rpcStream
	}

	go rpcStream.rpcClientRecv(fr, mqClose)
	return nil
}

// RPCStream 单个RPCStream
type RPCStream struct {
	GSID    string                  // 游戏服ID;e.g.: game1,game2
	Stream  pb.Service_StreamClient // 后端游戏服数据流
	IsClose bool
}

func (s *RPCStream) close() {
	if s.Stream != nil {
		err := s.Stream.CloseSend()
		if err != nil {
			log.Debugln("stream close err:", err)
		}
	}
}

// 读取rpc返回消息的goroutine
func (s *RPCStream) rpcClientRecv(mq chan pb.Frame, isMQClose chan struct{}) {
	log.Debugln("rpcClientRecv start")
	defer func() {
		log.Debugln("rpcClientRecv close")
	}()
	for {
		in, err := s.Stream.Recv()
		if err == io.EOF { // 流关闭
			log.Debugln("rpc stream client io eof", err)
			return
		}
		if err != nil {
			log.Debugln("rpc stream client err", err)
			return
		}
		select {
		case <-isMQClose: //多个生产者，当消费者触发信号时，直接返回
			return
		case mq <- *in:
		default:
			log.Debugln("mq chan full")
			return
		}
	}
}

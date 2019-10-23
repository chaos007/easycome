package grpc

import (
	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/enum"
)

var serverClient *Session

// GetServerClientSession 获得服务器的rpc客户端
func GetServerClientSession() *Session {
	if serverClient != nil && serverClient.Flag&enum.SessKickedOut == 0 {
		return serverClient
	}
	serverClient = NewSession()
	var stream pb.Service_StreamServer

	go serverClient.handle(stream)
	return serverClient
}

package grpc

import (
	"fmt"
	"github.com/chaos007/easycome/pb"
)

var serverClient *Session

// session常用量
const (
	SessKeyExchange = 0x1 // 是否已经交换完毕KEY
	SessEncrypt     = 0x2 // 是否可以开始加密
	SessKickedOut   = 0x4 // 踢掉
	SessAuthorized  = 0x8 // 已授权访问
	// SessHasClose    = 0x16 // 已授权访问
)

// GetServerClientSession 获得服务器的rpc客户端
func GetServerClientSession() *Session {
	if serverClient != nil && serverClient.Flag&SessKickedOut == 0 {
		return serverClient
	}
	serverClient = NewSession()
	var stream pb.Service_StreamServer

	if serverClient != nil {
		go func() {
			err := serverClient.handle(stream)
			if err != nil {
				fmt.Println("GetServerClientSession handle err:", err)
			}
		}()
	}
	return serverClient
}

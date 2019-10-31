package mixkcp

import (
	"github.com/chaos007/easycome/libs/msgmeta"
	"log"

	"github.com/golang/protobuf/proto"
)

var mapProtocol = map[int32]func(*Session, proto.Message) (proto.Message, error){}

// RegisterProtocol 注册协议
func RegisterProtocol(msgName string, m func(*Session, proto.Message) (proto.Message, error)) {
	v, ok := msgmeta.MessageMetaByName(msgName)
	if !ok {
		log.Fatal("wrong msgName")
		return
	}

	if _, ok := mapProtocol[v]; ok {
		log.Fatalf("duplicated protocol id:%d,name:%s", v, msgName)
		return
	}
	mapProtocol[v] = m
}

// getProtocol 获得协议函数
func getProtocol(id int32) func(*Session, proto.Message) (proto.Message, error) {
	if v, ok := mapProtocol[id]; ok {
		return v
	}
	return nil
}

package msgmeta

import (
	"github.com/chaos007/easycome/enum"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	msgNameToID     = make(map[string]int32)
	msgIDToName     = make(map[int32]string)
	msgToServerType = make(map[int32]string)
)

// RegisterMessageMeta 注册消息元信息(代码生成专用)
func RegisterMessageMeta(name string, id int32) bool {
	if _, ok := msgNameToID[name]; ok {
		log.Error("msg name repeated!")
		os.Exit(-1)
		return false
	}
	msgNameToID[name] = id
	if _, ok := msgIDToName[id]; ok {
		log.Error("msg id repeated!")
		os.Exit(-1)
		return false
	}
	msgIDToName[id] = name
	if strings.HasSuffix(name, enum.ToAgentString) {
		msgToServerType[id] = enum.ServerTypeAgent
	} else if strings.HasSuffix(name, enum.ToGameString) {
		msgToServerType[id] = enum.ServerTypeGame
	} else if strings.HasSuffix(name, enum.ToAllString) {
		msgToServerType[id] = enum.ServerTypeAll
	} else if strings.HasSuffix(name, enum.ToCenterString) {
		msgToServerType[id] = enum.ServerTypeCenter
	} else if strings.HasSuffix(name, enum.ToLoginString) {
		msgToServerType[id] = enum.ServerTypeLogin
	}
	return true
}

// GetMsgServerType 获得消息类型，转到相对应的服务器id
func GetMsgServerType(id int32) string {
	if v, ok := msgToServerType[id]; ok {
		return v
	}
	return ""
}

// MessageMetaByName 根据名字查找id
func MessageMetaByName(name string) (int32, bool) {
	if v, ok := msgNameToID[name]; ok {
		return v, true
	}

	return 0, false
}

// MessageMetaByID 根据id查找消息元信息
func MessageMetaByID(id int32) (string, bool) {
	if v, ok := msgIDToName[id]; ok {
		return v, true
	}

	return "", false
}

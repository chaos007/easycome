package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/chaos007/easycome/enum"

	log "github.com/sirupsen/logrus"
)

// General General
type General struct {
	ReadDeadLine int
	QueueLen     int
	SockBuf      int
	UDPSockBuf   int
	UDPSndWnd    int
	UDPRcvWnd    int
	UDPMtu       int
	Dscp         int
	Nodelay      int
	Interval     int
	ReSend       int
	Nc           int
	RPMlimit     int
}

// Config 服务器信息
type Config struct {
	General   General
	Center    CenterType
	Unique    CenterType
	AgentList []SessionType
	GameList  []SessionType
}

// CenterType CenterType
type CenterType struct {
	MysqlConfig     string
	LogLevel        string
	WebListen       string
	GrpcListen      string
	ETCDHost        []string
	ETCDRoot        string
	ETCDWatch       []string
	ETCDServicePath string
	ETCDConfigPath  string
}

// SessionType 对外的session
type SessionType struct {
	MysqlConfig     string
	ID              string
	MixListen       string //外部监听端口
	GrpcListen      string //rpc监听端口
	ProxyListen     string //代理服务器监听端口
	Version         string
	ETCDRoot        string   //etcd的服务器
	ETCDWatch       []string //etcd监听的目录
	ETCDServicePath string   //etcd自己所在的目录
	LogLevel        string   //log等级
	PProfListen     string   //pprof的端口
}

var dataConfig *Config

// ConfigName 系统配置文件
var ConfigName = "server.config"

// GetServerConfig 获得服务器配置
func GetServerConfig() *Config {
	return dataConfig
}

// GetServerMixkcpListen GetServerMixkcpListen
func GetServerMixkcpListen(key, severType string) string {
	if severType == enum.ServerTypeAgent {
		for _, item := range dataConfig.AgentList {
			if item.ID == key {
				return item.MixListen
			}
		}
	}

	if severType == enum.ServerTypeGame {
		for _, item := range dataConfig.GameList {
			if item.ID == key {
				return item.MixListen
			}
		}
	}

	return ""
}

// GetServerProxyListen GetServerMixkcpListen
func GetServerProxyListen(key, severType string) string {
	if severType == enum.ServerTypeAgent {
		for _, item := range dataConfig.AgentList {
			if item.ID == key {
				return item.ProxyListen
			}
		}
	}

	if severType == enum.ServerTypeGame {
		for _, item := range dataConfig.GameList {
			if item.ID == key {
				return item.ProxyListen
			}
		}
	}

	return ""
}

// ServerConfigInit 中心服务器配置表初始化
func ServerConfigInit() {
	dataConfig = new(Config)
	if _, err := os.Stat(ConfigName); os.IsNotExist(err) {
		log.Error("can not find server.config")
		os.Exit(-1)
		return
	}
	jsonByte, err := ioutil.ReadFile(ConfigName)
	if err != nil {
		log.Error("read server.config err", err.Error())
		os.Exit(-1)
		return
	}
	jsonByteReal := []byte{}
	for index := 0; index < len(jsonByte); index++ {
		if jsonByte[index] != '\n' && jsonByte[index] != ' ' && jsonByte[index] != '\t' {
			jsonByteReal = append(jsonByteReal, jsonByte[index])
		}
	}
	if err = json.Unmarshal(jsonByteReal, dataConfig); err != nil {
		log.Errorf("json Unmarshal server.config error:%s", err.Error())
		os.Exit(-1)
		return
	}
}

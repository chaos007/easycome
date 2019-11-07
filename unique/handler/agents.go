package handler

import (
	"errors"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/libs/enum"
	"github.com/chaos007/easycome/libs/config"
	"github.com/chaos007/easycome/libs/etcdservices"
	"github.com/chaos007/easycome/libs/hashs"
	"github.com/chaos007/easycome/model/account"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

//AgentInfo 服务器信息
type AgentInfo struct {
	ID     int
	Name   string
	AddrIP string
}

var idcount int
var list *AgentServers

//AgentServers 网关服务器列表
type AgentServers struct {
	Agents     map[int]*AgentInfo
	Name       string
	Consistent *hashs.Consistent
	lock       *sync.RWMutex
}

func (a *AgentServers) addAgent(info *AgentInfo) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, ok := a.Agents[info.ID]; ok {
		return false
	}
	a.Agents[info.ID] = info
	return true
}

func (a *AgentServers) delAgent(id int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	delete(a.Agents, id)
}

func (a *AgentServers) getAgentIDbyName(name string) int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for key, item := range a.Agents {
		if item.Name == name {
			return key
		}
	}
	return -1
}

func newAgentServer() *AgentServers {
	return &AgentServers{
		Agents:     map[int]*AgentInfo{},
		Consistent: hashs.NewConsistent(),
		lock:       new(sync.RWMutex),
	}
}

//Init init
func Init() {
	list = newAgentServer()
	agents, res := etcdservices.GetServiceDirRPCClient(enum.ServerTypeAgent)
	if res != true {
		log.Debugln("GetServiceDir no true")
	}
	for _, item := range agents {
		// groupName := filepath.Base(filepath.Dir(filepath.Dir(item.GetKey())))
		idcount++

		mixkcpListen := config.GetServerProxyListen(item.GetKey(), enum.ServerTypeAgent)
		log.Debugln("mixkcpListen:", mixkcpListen)
		list.addAgent(&AgentInfo{idcount, item.GetKey(), mixkcpListen})

	}
	etcdservices.AddWatchCallbacks(int32(clientv3.EventTypePut), PutServiceHandle)
	etcdservices.AddWatchCallbacks(int32(clientv3.EventTypeDelete), DeleteServiceHandle)
}

//PutServiceHandle 新增网关回调
func PutServiceHandle(etcdPath string, addr string) {
	log.Debugln("PutServiceHandle", etcdPath)

	// groupName := filepath.Base(filepath.Dir(filepath.Dir(etcdPath)))
	mixkcpListen := config.GetServerProxyListen(filepath.Base(etcdPath), enum.ServerTypeAgent)
	log.Debugln("PutServiceHandle mixkcpListen:", mixkcpListen)
	server := &AgentInfo{
		AddrIP: mixkcpListen,
		Name:   filepath.Base(etcdPath),
	}

	ipList := strings.Split(server.AddrIP, ":")

	if len(ipList) != 2 {
		log.Errorln("ip parse error,ip:", server.AddrIP)
	}
	if list == nil {
		list = newAgentServer()
	}
	idcount++
	server.ID = idcount
	list.addAgent(server)

	gaNode := hashs.NewNode(idcount, server.AddrIP, filepath.Base(etcdPath), 1)
	log.Debug("addGameServer ganode", gaNode)
	res := list.Consistent.Add(gaNode)
	if !res {
		log.Errorln("PutServiceHandle add node err,agentName:", filepath.Base(etcdPath))
	}
}

//DeleteServiceHandle 删除网关回调
func DeleteServiceHandle(etcdPath string, addr string) {
	log.Debugln("DeleteServiceHandle", etcdPath)
	// groupName := filepath.Base(filepath.Dir(filepath.Dir(etcdPath)))

	id := list.getAgentIDbyName(filepath.Base(etcdPath))
	list.delAgent(id)
	gaNode, exist := list.Consistent.GetNodeByID(id)
	if !exist {
		log.Errorln("DeleteServiceHandle del node err,agentName:", filepath.Base(etcdPath))
	}
	if list.Consistent != nil {
		list.Consistent.Remove(gaNode)
	}
}

func getAllGroupBalanceAgent(uid string) (*pb.WebDownGetServerInfo, error) {
	result := &pb.WebDownGetServerInfo{}
	gaNode, existsNode := list.Consistent.Get(uid)
	if !existsNode {
		return nil, nil
	}
	list := strings.Split(gaNode.Addr, ":")
	if len(list) != 2 {
		log.Errorln("can not parse ip addr")
		return result, errors.New("addr error")
	}
	result.Info = &pb.ServerInfo{
		AgentIP:   list[0],
		AgentPort: list[1],
	}
	return result, nil
}

// UserMap 用户当前的登陆时间
type UserMap struct {
	lock              *sync.RWMutex
	IDToLastLoginTime map[int64]*account.Account //
}

// CuidToUserID cuid对应UserID
type CuidToUserID struct {
	cuidToUserID map[string]map[string]*account.Account
	lock         *sync.RWMutex
}

// GetUserIDByCuid GetUserIDByCuid
func (c *CuidToUserID) GetUserIDByCuid(cuid, channel string) (*account.Account, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	v, ok := c.cuidToUserID[channel]
	if !ok {
		return nil, false
	}
	info, cok := v[cuid]

	return info, cok
}

// UpdateAcount UpdateAcount
func (c *CuidToUserID) UpdateAcount(a *account.Account) {
	// c.lock.Lock()
	// defer c.lock.Unlock()
	// if v, ok := c.cuidToUserID[a.Channel]; ok {
	// 	v[a.Cuid] = a
	// } else {
	// 	m := map[string]*account.Account{
	// 		a.Cuid: a,
	// 	}
	// 	c.cuidToUserID[a.Channel] = m
	// }
}

// UpdatePerson UpdatePerson
func (u *UserMap) UpdatePerson(a *account.Account) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.IDToLastLoginTime[a.ID] = a
}

// DelPerson DelPerson
func (u *UserMap) DelPerson(id int64) {
	u.lock.Lock()
	defer u.lock.Unlock()
	delete(u.IDToLastLoginTime, id)
}

// GetPerson GetPerson
func (u *UserMap) GetPerson(id int64) *account.Account {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.IDToLastLoginTime[id]
}

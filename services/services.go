package services

import (
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"path/filepath"
	"strings"
)

//服务器列表
var (
	_defaultPool = ServicePool{
		services: map[string]*service{},
	}
)

//RPCClient a single connection
type RPCClient struct {
	key  string
	conn *grpc.ClientConn
	addr string
}

// a kind of service
type service struct {
	clients map[string]*RPCClient
	idx     uint32 // for round-robin purpose
	mu      sync.RWMutex
}

// ServicePool all services
type ServicePool struct {
	root     string
	services map[string]*service
	mu       sync.RWMutex
}

// ServerUpdate 服务器更新
func ServerUpdate(m map[string]string, me string) {
	_defaultPool.mu.Lock()
	_defaultPool.services = map[string]*service{}
	_defaultPool.mu.Unlock()
	for k, v := range m {
		if k == me {
			continue
		}
		_defaultPool.addService(k, v)
	}
}

// add a service
func (p *ServicePool) addService(key, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// name check
	serviceName := strings.Replace(key, "\\", "/", -1) //防止不同平台文件地址字符出现问题
	serviceType := filepath.Base(filepath.Dir(serviceName))

	// try new service kind init
	if p.services[serviceType] == nil {
		p.services[serviceType] = &service{
			clients: map[string]*RPCClient{},
		}
	}

	// create service connection
	service := p.services[serviceType]
	if _, ok := service.clients[key]; ok {
		logrus.Debugln("master has already has service")
		return
	}
	logrus.Debugln("master addService begin key value", key, value)
	go func() {
		if conn, err := grpc.Dial(value, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
			service.mu.Lock()
			service.clients[key] = &RPCClient{key, conn, value}
			service.mu.Unlock()
			logrus.Debugln("service :", serviceType, " added:", key, "-->", value)
		} else {
			logrus.Debugln("did not connect:", key, "-->", value, "error:", err)
		}
	}()
}

// remove a service TODO
func (p *ServicePool) removeService(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// name check
	logrus.Debugln("master removeService key value", key)
	serviceName := strings.Replace(key, "\\", "/", -1) //防止不同平台文件地址字符出现问题
	serviceType := filepath.Base(filepath.Dir(serviceName))
	// check service kind
	service := p.services[serviceType]
	if service == nil {
		logrus.Debugln("no such service:", serviceName)
		return
	}

	// remove a service
	v, ok := service.clients[key]
	if !ok {
		logrus.Debugln("master did not have service")
		return
	}
	v.conn.Close()
	delete(service.clients, serviceName)
}

// getServiceWithID 根据id来获取服务
func (p *ServicePool) getServiceWithID(path string, id string) *grpc.ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// check existence
	service := p.services[path]
	if service == nil {
		return nil
	}
	if len(service.clients) == 0 {
		return nil
	}

	// loop find a service with id
	fullpath := path + "/" + id
	for k := range service.clients {
		if service.clients[k].key == fullpath {
			return service.clients[k].conn
		}
	}

	return nil
}

// get a service in round-robin style
// especially useful for load-balance with state-less services
func (p *ServicePool) getService(path string) (conn *grpc.ClientConn, key string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// check existence
	logrus.Debugln("getService path", path)
	service := p.services[path]
	if service == nil {
		return nil, ""
	}

	if len(service.clients) == 0 {
		return nil, ""
	}

	list := []*RPCClient{}

	for _, v := range service.clients {
		list = append(list, v)
	}

	// get a service in round-robind style,
	idx := int(atomic.AddUint32(&service.idx, 1)) % len(list)
	return service.clients[list[idx].key].conn, service.clients[list[idx].key].key
}

// GetAllService get all service conn
func (p *ServicePool) GetAllService() map[string]*grpc.ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// check existence
	logrus.Debugln("GetAllService services", p.services)
	result := map[string]*grpc.ClientConn{}
	for _, clients := range p.services {
		for _, c := range clients.clients {
			result[c.key] = c.conn
		}
	}
	return result
}

// GetServiceDirRPCClient get all service conn in the dir
func (p *ServicePool) GetServiceDirRPCClient(path string) map[string]*RPCClient {
	p.mu.RLock()
	defer p.mu.RUnlock()

	logrus.Debugln("GetServiceDirRPCClient path", path)
	clients, ok := p.services[path]
	if !ok {
		return nil
	}
	result := map[string]*RPCClient{}
	for _, c := range clients.clients {
		result[c.key] = c
	}
	return result
}

// GetService GetService
func GetService(path string) *grpc.ClientConn {
	conn, _ := _defaultPool.getService(path)
	return conn
}

// GetService2 GetService2
func GetService2(path string) (*grpc.ClientConn, string) {
	conn, key := _defaultPool.getService(path)
	return conn, key
}

// GetServiceWithID GetServiceWithID
func GetServiceWithID(path string, id string) *grpc.ClientConn {
	return _defaultPool.getServiceWithID(path, id)
}

// GetServiceDirRPCClient 获得目录的服务器
func GetServiceDirRPCClient(path string) map[string]*RPCClient {
	return _defaultPool.GetServiceDirRPCClient(path)
}

// GetAllService GetAllService
func GetAllService() map[string]*grpc.ClientConn {
	return _defaultPool.GetAllService()
}

// GetKey GetKey
func (r *RPCClient) GetKey() string {
	return r.key
}

// GetAddr GetAddr
func (r *RPCClient) GetAddr() string {
	return r.addr
}

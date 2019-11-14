package etcdservices

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
)

//RPCClient a single connection
type RPCClient struct {
	key  string
	conn *grpc.ClientConn
	addr string
}

// a kind of service
type service struct {
	clients map[string]RPCClient
	idx     uint32 // for round-robin purpose
}

// ServicePool all services
type ServicePool struct {
	root       string
	services   map[string]*service
	etcdClient *clientv3.Client
	mu         sync.RWMutex
	me         *Service
	watchFirst map[string]bool //每个监控列表第一次查看
}

//服务器列表
var (
	_defaultPool ServicePool
	once         sync.Once
)

//回调函数
var (
	_watchCallback WatchCallbackMgr
)

// Init 初始化
func Init(endpoints []string, root string, watchPath []string, me *Service) {
	once.Do(func() {
		_defaultPool.init(endpoints, root, watchPath, me)
	})
}

func (p *ServicePool) init(endpoints []string, root string, watchPath []string, me *Service) (*ServicePool, error) {
	// init etcd client
	log.Println("init Endpoints", endpoints)
	log.Println("init watchPath", watchPath)
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 1 * time.Second,
	}
	etcdcli, err := clientv3.New(cfg)
	if err != nil {
		log.Errorln("init etcd err:", err)
		os.Exit(-1)
		return nil, nil
	}

	p.etcdClient = etcdcli
	p.root = root
	p.services = make(map[string]*service)
	p.me = me
	p.watchFirst = make(map[string]bool)

	for _, path := range watchPath {
		go p.WatchNodes(me.servicePath+me.Name, path)
	}
	return p, err
}

// getServerOnce 获得当前etcd的服务器
func (p *ServicePool) getServerOnce(path string) {
	p.watchFirst[path] = true
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	resp, err := p.etcdClient.Get(ctx, path, clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Errorln("watch etcd err:", err)
		os.Exit(-1)
		return
	}

	if resp != nil && resp.Kvs != nil {
		log.Debugln("getServerOnce resp kvs length :", len(resp.Kvs))
		for _, item := range resp.Kvs {
			info := GetServiceInfoByValue(item.Value)
			if string(item.Key) == p.me.servicePath+p.me.Name {
				continue
			}
			p.addService(string(item.Key), info.IP)
			InvokeWatchCallBacks(int32(clientv3.EventTypePut), string(item.Key), info.IP)
		}
	}
}

// WatchNodes WatchNodes
func (p *ServicePool) WatchNodes(except, path string) {
	rch := p.etcdClient.Watch(context.TODO(), path, clientv3.WithPrefix())
	if v, ok := p.watchFirst[path]; !ok || !v {
		p.getServerOnce(path)
	}
	for {
		select {
		case wresp, ok := <-rch:
			if !ok {
				log.Errorln("etcd watch root die")
			} else {
				for _, ev := range wresp.Events {
					if string(ev.Kv.Key) == except {
						continue
					}
					switch ev.Type {
					case clientv3.EventTypePut:
						log.Debugf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
						info := GetServiceInfo(ev)
						p.addService(string(ev.Kv.Key), info.IP)
						InvokeWatchCallBacks(int32(clientv3.EventTypePut), string(ev.Kv.Key), info.IP)
					case clientv3.EventTypeDelete:
						log.Debugf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
						p.removeService(string(ev.Kv.Key))
						InvokeWatchCallBacks(int32(clientv3.EventTypeDelete), string(ev.Kv.Key), "")
					}
				}
			}
		}
	}
}

// add a service
func (p *ServicePool) addService(key, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.me.servicePath+p.me.Name == key {
		return
	}

	// name check
	serviceName := strings.Replace(key, "\\", "/", -1) //防止不同平台文件地址字符出现问题
	serviceType := filepath.Base(filepath.Dir(serviceName))

	// try new service kind init
	if p.services[serviceType] == nil {
		p.services[serviceType] = &service{
			clients: map[string]RPCClient{},
		}
	}

	// create service connection
	service := p.services[serviceType]
	if _, ok := service.clients[key]; ok {
		log.Debugln("master has already has service")
		return
	}
	log.Debugln("master addService begin key value", key, value)
	if conn, err := grpc.Dial(value, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
		service.clients[key] = RPCClient{key, conn, value}
		log.Println("service :", serviceType, " added:", key, "-->", value)
	} else {
		log.Errorln("did not connect:", key, "-->", value, "error:", err)
	}
}

// remove a service TODO
func (p *ServicePool) removeService(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// name check
	log.Debugln("master removeService key value", key)
	serviceName := strings.Replace(key, "\\", "/", -1) //防止不同平台文件地址字符出现问题
	serviceType := filepath.Base(filepath.Dir(serviceName))
	// check service kind
	service := p.services[serviceType]
	if service == nil {
		log.Debugln("no such service:", serviceName)
		return
	}

	// remove a service
	v, ok := service.clients[key]
	if !ok {
		log.Debugln("master did not have service")
		return
	}
	v.conn.Close()
	delete(service.clients, serviceName)
	return
}

// provide a specific key for a service, eg:
// path:/backends/snowflake, id:s1
//
// the full cannonical path for this service is:
// 			/backends/snowflake/s1
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
	fullpath := string(path) + "/" + id
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
	log.Debugln("getService path", path)
	service := p.services[path]
	if service == nil {
		return nil, ""
	}

	if len(service.clients) == 0 {
		return nil, ""
	}

	list := []RPCClient{}

	for _, v := range service.clients {
		list = append(list, v)
	}

	// get a service in round-robind style,
	idx := int(atomic.AddUint32(&service.idx, 1)) % len(list)
	return service.clients[list[idx].key].conn, service.clients[list[idx].key].key
}

// GetServiceDir get all service conn in the dir
func (p *ServicePool) GetServiceDir(path string) (map[string]*grpc.ClientConn, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// check existence
	log.Debugln("GetServiceDir path", path)
	log.Debugln("GetServiceDir services", p.services)
	clients, ok := p.services[path]
	if !ok {
		return nil, false
	}
	result := map[string]*grpc.ClientConn{}
	for _, c := range clients.clients {
		result[c.key] = c.conn
	}
	if len(result) <= 0 {
		return nil, false
	}
	return result, true
}

// GetAllService get all service conn
func (p *ServicePool) GetAllService() (map[string]*grpc.ClientConn, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// check existence
	log.Debugln("GetAllService services", p.services)
	result := map[string]*grpc.ClientConn{}
	for _, clients := range p.services {
		for _, c := range clients.clients {
			_, ok := result[c.key]
			if ok {
				return nil, false
			}
			result[c.key] = c.conn
		}
	}
	if len(result) <= 0 {
		return nil, false
	}
	return result, true
}

// GetServiceDirRPCClient get all service conn in the dir
func (p *ServicePool) GetServiceDirRPCClient(path string) (map[string]*RPCClient, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// check existence
	log.Debugln("GetServiceDirRPCClient path", path)
	clients, ok := p.services[path]
	if !ok {
		return nil, false
	}
	result := map[string]*RPCClient{}
	for _, c := range clients.clients {
		_, ok := result[c.key]
		if ok {
			return nil, false
		}
		result[c.key] = &c
	}
	if len(result) <= 0 {
		return nil, false
	}
	return result, true
}

// GetService GetService
func GetService(path string) *grpc.ClientConn {
	conn, _ := _defaultPool.getService(path)
	return conn
}

// GetService2 GetService2
func GetService2(path string) (*grpc.ClientConn, string) {
	conn, key := _defaultPool.getService(_defaultPool.root + "/" + path)
	return conn, key
}

// GetServiceWithID GetServiceWithID
func GetServiceWithID(path string, id string) *grpc.ClientConn {
	return _defaultPool.getServiceWithID(_defaultPool.root+"/"+path, id)
}

// GetServiceDirRPCClient 获得目录的服务器
func GetServiceDirRPCClient(path string) (map[string]*RPCClient, bool) {
	return _defaultPool.GetServiceDirRPCClient(path)
}

// GetServiceDir 获得目录的服务器
func GetServiceDir(path string) (map[string]*grpc.ClientConn, bool) {
	return _defaultPool.GetServiceDir(path)
}

// GetAllService GetAllService
func GetAllService() (map[string]*grpc.ClientConn, bool) {
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

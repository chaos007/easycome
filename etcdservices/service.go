package etcdservices

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/chaos007/easycome/config"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

// ServiceInfo the detail of service
type ServiceInfo struct {
	IP string
}

var me *Service

// GetMe 获得自己
func GetMe() *Service {
	return me
}

// Service service
type Service struct {
	Name        string
	Info        ServiceInfo
	stop        chan error
	leaseid     clientv3.LeaseID
	client      *clientv3.Client
	General     config.General
	servicePath string
	infoPlace   string
	serviceType string
	group       string
	version     string
}

// SetService 设置server的信息
func (s *Service) SetService(name string, info ServiceInfo, servicePath string) {
	s.Name = name
	s.Info = info
	s.stop = make(chan error)
	s.servicePath = servicePath
}

// Start Start
func (s *Service) Start() error {

	ch, err := s.keepAlive()
	if err != nil {
		log.Fatalln("etcd keepAlive err:", err)
		return err
	}

	for {
		select {
		case err := <-s.stop:
			s.revoke()
			return err
		case <-s.client.Ctx().Done():
			log.Fatalln("server has closed")
			return errors.New("server closed")
		case _, ok := <-ch:
			if !ok {
				log.Debugln("keep alive channel closed")
				s.revoke()
				return nil
			}
			// log.Debugf("Recv reply from service: %s, ttl:%d", s.Name, ka.TTL)
		}
	}
}

// Stop Stop
func (s *Service) Stop() {
	s.stop <- nil
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	info := &s.Info

	key := s.servicePath + s.Name
	value, _ := json.Marshal(info)

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Debugln("keepAlive put key", key, string(value))
	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s.leaseid = resp.ID

	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Service) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("servide:%s stop\n", s.Name)
	return err
}

// WatchNodes WatchNodes
func (s *Service) WatchNodes(path string, callback func()) {
	rch := s.client.Watch(context.TODO(), path, clientv3.WithPrefix())
	log.Debugln("event watchPath len", len(rch))
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				log.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				callback()
			}
		}
	}
}

// GetServiceInfo GetServiceInfo
func GetServiceInfo(ev *clientv3.Event) *ServiceInfo {
	info := &ServiceInfo{}
	err := json.Unmarshal([]byte(ev.Kv.Value), info)
	if err != nil {
		return nil
	}
	return info
}

// GetServiceInfoByValue GetServiceInfoByValue
func GetServiceInfoByValue(value []byte) *ServiceInfo {
	info := &ServiceInfo{}
	err := json.Unmarshal(value, info)
	if err != nil {
		return nil
	}
	return info
}

// GetClient etcd的客户端
func (s *Service) GetClient() *clientv3.Client {
	return s.client
}

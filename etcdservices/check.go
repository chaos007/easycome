package etcdservices

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/chaos007/easycome/enum"
	"github.com/chaos007/easycome/config"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

// BeforeService 请求配置数据前的客户端
func BeforeService(endpoints []string) (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 1 * time.Second,
	})

	if err != nil {
		return nil, err
	}
	s := &Service{
		client: cli,
	}

	me = s
	return s, err
}

// PutKey 放置值
func (s *Service) PutKey(key string, value string, ttl int64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(ttl)*time.Second)
	defer cancel()
	_, err := s.client.Put(ctx, key, value)
	if err != nil {
		logrus.Errorln("etcd PutKey error:", err.Error())
		return err
	}

	return nil
}

// GetValue 获取值
func (s *Service) GetValue(key string, ttl int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(ttl)*time.Second)
	defer cancel()
	r, err := s.client.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if len(r.Kvs) != 1 {
		return "", errors.New("wrong value length")
	}

	return string(r.Kvs[0].Value), nil
}

// CheckMe 跟etcd检查自身
func (s *Service) CheckMe(serverInfo, name, serverType string) (config.SessionType, bool) {
	sessionType := config.SessionType{}
	serverInfoString, err := s.GetValue(serverInfo, 10)
	if err != nil {
		logrus.Errorln("etcd GetValue ,err:", err)
		return sessionType, false
	}

	// 判断区组服务器，最多一组有100个
	// if !uuid.InitIndex(name) {
	// 	return nil, false
	// }
	si := config.Config{}

	s.infoPlace = serverInfo
	s.Name = name
	s.serviceType = serverType

	err = json.Unmarshal([]byte(serverInfoString), &si)
	if err != nil {
		logrus.Errorln("can not marshal server info ,err:", err)
		return sessionType, false
	}

	s.General = si.General

	if serverType == enum.ServerTypeAgent {
		for _, item := range si.AgentList {
			if item.ID == name {
				s.version = item.Version
				return item, true
			}
		}
	}

	if serverType == enum.ServerTypeGame {
		for _, item := range si.GameList {
			if item.ID == name {
				s.version = item.Version
				return item, true
			}
		}
	}

	return sessionType, false
}

// WatchMeChange 监控自己修改
func (s *Service) WatchMeChange(ch chan int) {
	for {
		rch := s.client.Watch(context.TODO(), s.infoPlace, clientv3.WithPrefix())
		logrus.Debugln("WatchMeChange event watchPath len", len(rch))
		select {
		case wresp, ok := <-rch:
			if !ok {
				logrus.Errorln("etcd watch root die")
			} else {
				for _, ev := range wresp.Events {
					si := config.Config{}

					logrus.Debugln("watch me change happen")
					err := json.Unmarshal([]byte(ev.Kv.Value), &si)
					if err != nil {
						logrus.Errorln("can not marshal server info ,err:", err)
					}

					if s.serviceType == enum.ServerTypeAgent {
						var hasMe bool
						for _, item := range si.AgentList {
							if item.ID == s.Name {
								hasMe = true
								if item.Version != s.version {
									logrus.Debugln("agent version change")
									close(ch)
									return
								}
							}
						}
						if !hasMe {
							logrus.Debugln("agent me delete")
							close(ch)
							return
						}
					}

					if s.serviceType == enum.ServerTypeGame {
						var hasMe bool
						for _, item := range si.GameList {
							if item.ID == s.Name {
								hasMe = true
								if item.Version != s.version {
									logrus.Debugln("game version change")
									close(ch)
									return
								}
							}
						}
						if !hasMe {
							logrus.Debugln("game me delete")
							close(ch)
							return
						}
					}
				}
			}
		}
	}
}

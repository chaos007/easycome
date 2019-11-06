package mixkcp

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// methodList 函数列表
type methodList struct {
	List []func(string) error
	lock *sync.RWMutex
}

func newMethodList() *methodList {
	return &methodList{
		List: []func(string) error{},
		lock: &sync.RWMutex{},
	}
}

func (m *methodList) addMethod(method func(string) error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.List = append(m.List, method)
	log.Debugln("methodList addMethod:", len(m.List))
}

func (m *methodList) exce(userid string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, item := range m.List {
		item(userid)
	}
}

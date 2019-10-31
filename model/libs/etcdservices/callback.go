package etcdservices

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

//WatchCallbackMgr 回调管理
type WatchCallbackMgr struct {
	list            map[int32][]func(string, string)
	watchEventMutex sync.RWMutex
}

var watchCallbackMgr = &WatchCallbackMgr{
	list:            map[int32][]func(string, string){},
	watchEventMutex: sync.RWMutex{},
}

//AddWatchCallbacks ..
func AddWatchCallbacks(eventType int32, callback func(etcdPath string, addr string)) {
	watchCallbackMgr.addWatchCallbacks(eventType, callback)
}

//AddWatchCallbacks 添加回调函数
func (w *WatchCallbackMgr) addWatchCallbacks(eventType int32, callback func(etcdPath string, addr string)) {
	w.watchEventMutex.Lock()
	defer w.watchEventMutex.Unlock()
	if v, ok := watchCallbackMgr.list[eventType]; ok {
		v = append(v, callback)
	} else {
		list := []func(string, string){}
		list = append(list, callback)
		watchCallbackMgr.list[eventType] = list
	}
}

//InvokeWatchCallBacks ..
func InvokeWatchCallBacks(eventType int32, etcdPath string, addr string) {
	watchCallbackMgr.invokeWatchCallBacks(eventType, etcdPath, addr)
}

//invokeWatchCallBacks 执行回调函数
func (w *WatchCallbackMgr) invokeWatchCallBacks(eventType int32, etcdPath string, addr string) {
	w.watchEventMutex.Lock()
	defer w.watchEventMutex.Unlock()
	log.Debugln("invokeWatchCallBacks", eventType, etcdPath, addr)
	if v, ok := watchCallbackMgr.list[eventType]; ok {
		for _, item := range v {
			item(etcdPath, addr)
		}
	}
}

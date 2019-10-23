package grpc

import (
	"sync"
)

// Registry 注册用户
type Registry struct {
	records map[int64]interface{} // id -> v
	sync.RWMutex
}

var (
	_defaultRegistry Registry
)

func init() {
	_defaultRegistry.init()
}

func (r *Registry) init() {
	r.records = make(map[int64]interface{})
}

//Register 注册 user
func (r *Registry) Register(id int64, v interface{}) {
	r.Lock()
	r.records[id] = v
	r.Unlock()
}

//Unregister 注销 user
func (r *Registry) Unregister(id int64, v interface{}) {
	r.Lock()
	if oldv, ok := r.records[id]; ok {
		if oldv == v {
			delete(r.records, id)
		}
	}
	r.Unlock()
}

//Query 查询 user
func (r *Registry) Query(id int64) (x interface{}) {
	r.RLock()
	x = r.records[id]
	r.RUnlock()
	return
}

//Count 在线人数
func (r *Registry) Count() (count int) {
	r.RLock()
	count = len(r.records)
	r.RUnlock()
	return
}

// Register ...
func Register(id int64, v interface{}) {
	_defaultRegistry.Register(id, v)
}

// Unregister ...
func Unregister(id int64, v interface{}) {
	_defaultRegistry.Unregister(id, v)
}

// Query ...
func Query(id int64) interface{} {
	return _defaultRegistry.Query(id)
}

// Count ...
func Count() int {
	return _defaultRegistry.Count()
}

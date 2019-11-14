package data

import "sync"

// ToSync 需要同步的数据
type ToSync struct {
	M    Syncer
	lock *sync.RWMutex
}

var data = &ToSync{
	lock: new(sync.RWMutex),
}

// Update 更新最新数据
func (t *ToSync) Update(s Syncer) {

}

// Cover 覆盖
func (t *ToSync) Cover(s Syncer) {

}

// Syncer 需要同步数据的接口
type Syncer interface {
	Next() interface{}
	Name() string
	GetMap() sync.Map
}

// ItemSyncer 单项同步数据的接口
type ItemSyncer interface {
	Next()
	Left() int
}

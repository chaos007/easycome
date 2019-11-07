package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

// DesignConfig ...
type DesignConfig struct {
	lock *sync.RWMutex
}

var design = new(sync.Map)

var reflectTypeToName = map[reflect.Type]string{
	reflect.TypeOf((*ItemMap)(nil)): "item",
}

var nameToReflectType = map[string]reflect.Type{
	"item": reflect.TypeOf((*ItemMap)(nil)),
}

// RegisterType 注册类型
func RegisterType(name string, rType reflect.Type) {
	reflectTypeToName[rType] = name
	nameToReflectType[name] = rType
}

// Maper 配置表model实现的解析
type Maper interface {
	ForEach(callback func(Parser) interface{}) (r []interface{})
	ParseTalbe(str string) error
	SetItem(Parser)
	GetItem(id int64) Parser
	DelItem(id int64)
}

// Parser 配置表单项实现的
type Parser interface {
	ParserJSON(str string) error
	GenJSON() (string, error)
	GetID() int64
}

// GetDesignConfig 配置表信息
func GetDesignConfig() *sync.Map {
	return design
}

// DesignConfigInitByEtcd 配置表初始化，配置表放置etcd上
func DesignConfigInitByEtcd(client *clientv3.Client, path string, etcdPath string) (mList []Maper, err error) {
	design = new(sync.Map)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if !strings.HasSuffix(etcdPath, "/") {
		etcdPath += "/"
	}
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("read data config dir err:%s", err.Error())
		return nil, err
	}
	for _, item := range fileInfos {
		name := strings.TrimSuffix(item.Name(), filepath.Ext(item.Name()))
		jsonByte, err := ioutil.ReadFile(path + item.Name())
		if err != nil {
			log.Errorf("read name :%s ,err:%s", item.Name(), err.Error())
			return nil, err
		}
		v, ok := nameToReflectType[name]
		if !ok {
			log.Errorln("can not find struct:", item.Name())
			return nil, errors.New("can not find struct")
		}
		jsonByteReal := []byte{}
		for index := 0; index < len(jsonByte); index++ {
			if jsonByte[index] != '\n' && jsonByte[index] != ' ' && jsonByte[index] != '\t' {
				jsonByteReal = append(jsonByteReal, jsonByte[index])
			}
		}
		maper, ok := reflect.New(v.Elem()).Interface().(Maper)
		if err := maper.ParseTalbe(string(jsonByteReal)); err != nil {
			return nil, err
		}
		maper.ForEach(func(p Parser) (in interface{}) {
			s, err := p.GenJSON()
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			_, err = client.Put(ctx, etcdPath+name+"/"+strconv.FormatInt(p.GetID(), 10), s)
			if err != nil {
				cancel()
				return err
			}
			cancel()
			return
		})
		design.Store(name, maper)
		mList = append(mList, maper)
	}
	return mList, nil
}

// DesignConfigInit 配置表初始化，配置表放在中心服上
func DesignConfigInit(path string) (mList []Maper, err error) {
	design = new(sync.Map)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("read data config dir err:%s", err.Error())
		return nil, err
	}
	for _, item := range fileInfos {
		name := strings.TrimSuffix(item.Name(), filepath.Ext(item.Name()))
		jsonByte, err := ioutil.ReadFile(path + item.Name())
		if err != nil {
			log.Errorf("read name :%s ,err:%s", item.Name(), err.Error())
			return nil, err
		}
		v, ok := nameToReflectType[name]
		if !ok {
			log.Errorln("can not find struct:", item.Name())
			return nil, errors.New("can not find struct")
		}
		jsonByteReal := []byte{}
		for index := 0; index < len(jsonByte); index++ {
			if jsonByte[index] != '\n' && jsonByte[index] != ' ' && jsonByte[index] != '\t' {
				jsonByteReal = append(jsonByteReal, jsonByte[index])
			}
		}
		maper, ok := reflect.New(v.Elem()).Interface().(Maper)
		if err := maper.ParseTalbe(string(jsonByteReal)); err != nil {
			return nil, err
		}
		design.Store(name, maper)
		mList = append(mList, maper)
	}
	return mList, nil
}

// UpdateDesignConfig UpdateDesignConfig
func UpdateDesignConfig(data map[string]string) error {
	design = new(sync.Map)
	for name, item := range data {
		v, ok := nameToReflectType[name]
		if !ok {
			log.Errorln("can not find struct:", name)
			return errors.New("can not find struct")
		}
		maper, ok := reflect.New(v.Elem()).Interface().(Maper)
		if err := maper.ParseTalbe(item); err != nil {
			return err
		}
		design.Store(name, maper)
	}
	return nil
}

// UpdateDesignConfigByEtcd UpdateDesignConfigByEtcd
func UpdateDesignConfigByEtcd(path string, data string, isDelete bool) {
	idString := filepath.Base(path)
	dirString := filepath.Base(filepath.Dir(path))
	m, ok := design.Load(dirString)
	if !ok {
		m := new(Maper)
		design.Store(dirString, m)
	}
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return
	}
	if isDelete {
		m.(Maper).DelItem(id)
	}
	p := m.(Maper).GetItem(id)
	if p == nil { //新数据
		v, ok := nameToReflectType[dirString]
		if !ok {
			return
		}
		p = reflect.New(v.Elem()).Interface().(Parser)
		p.ParserJSON(data)
	} else {
		p.ParserJSON(data)
	}
}

// Watch 监控数据变化
func Watch(client *clientv3.Client, path string) {
	rch := client.Watch(context.TODO(), path, clientv3.WithPrefix())
	log.Debugln("event watchPath len", len(rch))
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				log.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				UpdateDesignConfigByEtcd(string(ev.Kv.Key), string(ev.Kv.Value), true)
			case clientv3.EventTypeDelete:
				log.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				UpdateDesignConfigByEtcd(string(ev.Kv.Key), string(ev.Kv.Value), false)
			}
		}
	}
}

// TestMap TestMap
type TestMap struct {
	M    map[int64]*Test
	lock sync.RWMutex
}

// UpdateData 更新数据
func (t *TestMap) UpdateData(id int64, str string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if v, ok := t.M[id]; ok {
		v.JSONString = str
		v.ParserJSON()
	} else {
		test := &Test{
			JSONString: str,
		}
		test.ParserJSON()
	}
}

// ForEach ForEach
func (t *TestMap) ForEach(callback func(*Test) interface{}) []interface{} {
	t.lock.Lock()
	defer t.lock.Unlock()
	result := []interface{}{}
	for _, item := range t.M {
		i := callback(item)
		result = append(result, i)
	}
	return result
}

// Init Init
func (t *TestMap) Init() {
	fmt.Println("---------33333333333")
}

// Test ..
type Test struct {
	ID         int64  `json:"id"`
	JSONString string `json:"-"`
}

// ParserJSON 解析自己的json
func (t *Test) ParserJSON() error {
	return json.Unmarshal([]byte(t.JSONString), &t)
}

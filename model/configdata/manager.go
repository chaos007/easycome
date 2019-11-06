package configdata

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

// Design 策划配置表
type Design struct {
	LevelData map[int64]*Level `json:"level"`
	ItemData  map[int64]*Item  `json:"item"`

	lock *sync.RWMutex
}

var design = &Design{
	lock: new(sync.RWMutex),
}

// GetDesign 获得配置表
func GetDesign() *Design {
	return design
}

var reflectTypeToName = map[reflect.Type]string{}

var nameToReflectType = map[string]reflect.Type{}

// registerType 注册类型
func registerType(name string, rType reflect.Type) {
	reflectTypeToName[rType] = name
	nameToReflectType[name] = rType
}

// UpdateDesignConfig UpdateDesignConfig
func UpdateDesignConfig(data string) error {
	design = &Design{
		lock: new(sync.RWMutex),
	}

	err := json.Unmarshal([]byte(data), &design)

	return err
}

// DesignConfigInit 配置表初始化，配置表放在中心服上
func DesignConfigInit(path string) (s string, err error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	s = "{"
	for index, item := range fileInfos {
		name := strings.TrimSuffix(item.Name(), filepath.Ext(item.Name()))
		jsonByte, err := ioutil.ReadFile(path + item.Name())
		if err != nil {
			return "", err
		}

		jsonByteReal := []byte{}
		for index := 0; index < len(jsonByte); index++ {
			if jsonByte[index] != '\n' && jsonByte[index] != ' ' && jsonByte[index] != '\t' {
				jsonByteReal = append(jsonByteReal, jsonByte[index])
			}
		}
		s += `"` + name + `":` + string(jsonByteReal)
		if index != len(fileInfos)-1 {
			s += ","
		}
	}
	s += "}"

	return s, UpdateDesignConfig(s)
}

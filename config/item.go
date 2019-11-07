package config

import (
	"encoding/json"
	"sync"
)

// ItemMap 整张表
type ItemMap struct {
	m    map[int64]*Item
	lock sync.RWMutex
}

// Item 道具表
type Item struct {
	ID          int64  `json:"id"`
	Type        int32  `json:"type"`
	Name        string `json:"name"`
	ItemSkillID int64  `json:"item_skill_id"`
	JSONString  string `json:"-"`
}

// ParseTalbe 解析一张json表
func (t *ItemMap) ParseTalbe(str string) error {
	itemMap := map[int64]*Item{}
	err := json.Unmarshal([]byte(str), &itemMap)
	if err != nil {
		return err
	}
	t.m = itemMap
	return nil
}

// ForEach ForEach
func (t *ItemMap) ForEach(callback func(Parser) interface{}) (r []interface{}) {
	t.lock.Lock()
	defer t.lock.Unlock()
	for _, item := range t.m {
		i := callback(item)
		r = append(r, i)
	}
	return
}

// GetItem 获得单项
func (t *ItemMap) GetItem(id int64) Parser {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.m[id]
}

// SetItem 设置单项
func (t *ItemMap) SetItem(p Parser) {
	t.lock.Lock()
	defer t.lock.Unlock()
	i := p.(*Item)
	t.m[i.ID] = i
}

// DelItem 删除单项
func (t *ItemMap) DelItem(id int64) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.m, id)
}

// ParserJSON 解析自己的json
func (t *Item) ParserJSON(str string) error {
	t.JSONString = str
	return json.Unmarshal([]byte(t.JSONString), &t)
}

// GenJSON 生成json
func (t *Item) GenJSON() (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	t.JSONString = string(b)
	return string(b), nil
}

// GetID 获得id
func (t *Item) GetID() int64 {
	return t.ID
}

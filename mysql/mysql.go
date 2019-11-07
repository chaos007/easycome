package mysql

import (
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

// GetEngine 后续查看是否需要加入池
func GetEngine() *xorm.Engine {
	return engine
}

// NewMySQL 新的sql
func NewMySQL(mysqlLink string) (err error) {
	engine, err = xorm.NewEngine("mysql", mysqlLink)
	if err != nil {
		return
	}
	engine.SetMaxOpenConns(2000)
	engine.SetMaxIdleConns(1000)
	return
}

// EngineSync2 同步表
func EngineSync2(beans interface{}) error {
	return engine.Sync2(beans)
}

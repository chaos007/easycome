package main

import (
	"github.com/chaos007/easycome/libs/mysql"
	"github.com/chaos007/easycome/model/player"

	"github.com/sirupsen/logrus"
)

// dbInit 连接数据库库
func dbInit(mysqlLink string) {
	if err := mysql.NewMySQL(mysqlLink); err != nil {
		logrus.Fatalf("mysql link:%s error:%s", mysqlLink, err.Error())
	}
	logrus.Println("mysql start ...")

	if err := mysql.EngineSync2(&player.Player{}); err != nil {
		logrus.Fatal("sync mysql err:", err.Error())
	}
}


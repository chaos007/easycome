package main

import (
	_ "github.com/go-sql-driver/mysql"
)

// tabelInit 连接库
func tabelInit() {
	// if err := mysql.NewMySQL(config.GetServerConfig().Center.MysqlConfig); err != nil {
	// 	logrus.Fatalf("mysql link:%s error:%s", config.GetServerConfig().Center.MysqlConfig, err.Error())
	// }
	// logrus.Println("mysql start ...")

	// if err := mysql.EngineSync2(&account.Account{}); err != nil {
	// 	logrus.Fatal("sync mysql err:", err.Error())
	// }
}

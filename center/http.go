package main

import (
	"github.com/chaos007/easycome/center/handle"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HTTPInit http服务初始化
func HTTPInit() {

	handle.HTTPInit()

	fileHandler := http.FileServer(http.Dir("./static"))
	http.Handle("/", ourLoggingHandler(fileHandler))

	// err := http.ListenAndServeTLS(":443", config.GetServerCongfig().CRTFile, config.GetServerCongfig().KeyFile, nil)
	//err := http.ListenAndServe(config.GetServerConfig().BackListen, nil) //设置监听的端口
	err := http.ListenAndServe(":8888", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ourLoggingHandler(h http.Handler) http.Handler {
	result := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		h.ServeHTTP(w, r)
	})
	return result
}

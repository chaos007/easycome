package main

import (
	"flag"
	"fmt"

	"github.com/chaos007/easycome/tools/testscript/client"
	"github.com/chaos007/easycome/tools/testscript/player"

	log "github.com/sirupsen/logrus"
)

//"github.com/chaos007/easycome/data/pb"

// "github.com/chaos007/easycome/tool/testscript/http"

//外网ip 42.62.64.254
//内网ip 192.168.51.151
//内网测试ip 192.168.62.207
//http://webapi.xyzcs.aiwan4399.com/

// GetAgentData 获得agent的数据结构
type GetAgentData struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data *IPPort `json:"data"`
}

// IPPort ip端口
type IPPort struct {
	IP        string
	Port      string
	BroadText string
}

var host = flag.String("h", "http://192.168.62.207", "please input server host") // *addr
var port = flag.String("p", "8066", "please input server port")
var num = flag.Int("num", 1, "please input client num")                                             //*num
var mac = flag.String("mac", "lantiao", "please input client mac")                                  //*mac
var debug = flag.Bool("debug", false, "need show debug log on screen?")                             //*debug
var battletype = flag.Int("bt", 3, "battle type,normal 1v1,3v3,5v5 =0,1,2,ROBOT 1v1,3v3,5v5=3,4,5") //*battletype
var idbegin = flag.Int("ib", 400001, "the user id start from")                                      //*idbegin
var testtype = flag.Int("tt", 1, "test type press test=1,black box test=2")                         //*idbegin

var die = make(chan bool)

var step = make(chan int)

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		arr := flag.Args()
		switch flag.NArg() {
		case 1:
		case 2:
		default:
		}
		fmt.Println(arr)
	}
	println("idbegin:", *idbegin)

	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Infof("host:%s", addr)
	if !*debug {
		log.SetLevel(log.ErrorLevel)
	}
	switch *testtype {
	case 1:
		for i := *idbegin; i < *idbegin+*num; i++ {
			Client := client.NewClient()
			Client.Session.UserID = int64(i)
			p := player.GetPlayer(int64(i))
			p.ID = int64(i)
			p.SetSession(Client.Session)

			// result, err := http.HttpPost(addr+"/getAgent", getAgentValue)
			// if err != nil {
			// 	println("err:", err.Error())
			// 	return
			// }
			// getAgentData := new(GetAgentData)
			// getAgentData.Data = new(IPPort)
			// json.Unmarshal([]byte(result), getAgentData)
			// fmt.Printf("ipport:%#v", *getAgentData)

			// useActiveValue := map[string][]string{
			// 	"code": []string{"X7vOLk"},
			// 	"uid":  []string{Client.Session.UserID},
			// 	"syn":  []string{http.Getmd5(Client.Session.UserID, "X7vOLk")},
			// }
			// resultUseActive, err := http.HttpPost(addr+"/useActive", useActiveValue)
			// if err != nil {
			// 	println("err:", err.Error())
			// 	return
			// }
			// activationData := new(http.JsonRet)
			// err = json.Unmarshal([]byte(resultUseActive), activationData)
			// if err != nil {
			// 	println("Activation check json error:", err.Error())
			// 	return
			// }
			// if activationData.Code != http.STATUS_OK {
			// 	return
			// }

			// go Client.PressTestUnit(getAgentData.Data.IP+":"+getAgentData.Data.Port, die, p, battletype)
			// go Client.PressTestUnit(*host+":"+*port, die, p, battletype)
		}
	}

	//Client.FuncTestList(addr, die)

	select {}
}

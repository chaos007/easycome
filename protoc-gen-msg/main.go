package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/alecthomas/template"
	"github.com/golang/protobuf/proto"

	pbcompiler "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var severTmplPath = "api_server.tmpl"
var clientTmplPath = "api_client.tmpl"
var targetMsgServerFile = "pb/msgRegister.go"
var targetMsgClientFile = "client/msgRegister.ts"

type apiExpr struct {
	Name string
	ID   int
}

func main() {

	dir := strings.TrimSuffix(targetMsgServerFile, path.Base(targetMsgServerFile))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal("target server dir not exist ", targetMsgServerFile)
		os.Exit(-1)
	}

	dir = strings.TrimSuffix(targetMsgClientFile, path.Base(targetMsgClientFile))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal("target client dir not exist ", targetMsgClientFile)
		os.Exit(-1)
	}

	if _, err := os.Stat(targetMsgServerFile); err == nil {
		os.Remove(targetMsgServerFile)
	}

	if _, err := os.Stat(targetMsgClientFile); err == nil {
		os.Remove(targetMsgClientFile)
	}

	pluginInput, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// 转换输入
	var req pbcompiler.CodeGeneratorRequest

	if err := proto.Unmarshal(pluginInput, &req); err != nil {
		fmt.Printf(err.Error())
		return
	}

	exprList := []apiExpr{}
	id := 0
	for i := 0; i < len(req.ProtoFile); i++ {
		messageType := (*req.ProtoFile[i]).MessageType
		for j := 0; j < len(messageType); j++ {
			id++
			expr := apiExpr{
				Name: messageType[j].GetName(),
				ID:   id,
			}

			exprList = append(exprList, expr)
		}
	}

	serverRegisterTmpl, err := template.New("api_server.tmpl").ParseFiles(severTmplPath)
	if err != nil {
		log.Fatal(err)
	}
	apiArgs := struct {
		Exprs []apiExpr
	}{exprList}

	var serverRegisterFile *os.File
	defer serverRegisterFile.Close()

	serverRegisterFile, err = os.OpenFile(targetMsgServerFile, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return
	}
	err = serverRegisterTmpl.Execute(serverRegisterFile, apiArgs)
	if err != nil {
		log.Fatal(err)
	}

	clientTmpl, err := template.New("api_client.tmpl").ParseFiles(clientTmplPath)
	if err != nil {
		log.Fatal(err)
	}

	var clientRegisterFile *os.File
	defer serverRegisterFile.Close()

	clientRegisterFile, err = os.OpenFile(targetMsgClientFile, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return
	}
	err = clientTmpl.Execute(clientRegisterFile, apiArgs)
	if err != nil {
		log.Fatal(err)
	}
}

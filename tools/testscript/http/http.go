package http

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
)

const (
	AgentKey = "|x-2w8329(#Dm)Q|"
)

const (
	UPDATE_ERROR = -1002 //   = "activation update error"
	SYN_ERROR    = -1003 //   = "activation syn error"
	STATUS_OK    = 0     //   = "activation status ok"
)

var retData = map[int]JsonRet{
	UPDATE_ERROR: JsonRet{Code: UPDATE_ERROR, Msg: "网络不畅,稍候再试"},
	SYN_ERROR:    JsonRet{Code: SYN_ERROR, Msg: "参数有误"},
	STATUS_OK:    JsonRet{Code: STATUS_OK, Msg: "成功"},
}

// JsonRet 返回给客户端的通用结构
type JsonRet struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"-"`
}

func HttpPost(path string, value map[string][]string) (string, error) {
	var err error
	resp, err := http.PostForm(path, value)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func Getmd5(name1, name2 string) string {
	var key string
	key = name1 + AgentKey + name2
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(key))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

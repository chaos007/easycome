package web

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// 错误编码
const (
	PasswordError = -10002
	SynError      = -10001
	ParamError    = -10000
	StatusOk      = 10000
)

// JSONRet 返回给客户端的通用结构
type JSONRet struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Key 通信密钥
var Key = "a40ef1cbe609853bc5f1d6c83caaca75"

// HTTPReturnWrite 数据返回
func HTTPReturnWrite(w http.ResponseWriter, r *JSONRet) {
	b, err := json.Marshal(r)
	if err != nil {
		logrus.Errorf("HTTPReturnWrite write json data err:%s", err.Error())
	}
	w.Write(b)
}

// RegularCheck 例行检查
func RegularCheck(w http.ResponseWriter, r *http.Request) (string, bool) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	var err error
	err = r.ParseForm()
	if err != nil {
		HTTPReturnWrite(w, &JSONRet{Code: SynError, Msg: "参数错误"})
		return "", false
	}
	data := Parse(r, "data")
	syn := Parse(r, "syn")
	if Getmd5(data, Key) != syn {
		HTTPReturnWrite(w, &JSONRet{Code: SynError, Msg: "校对码错误"})
		return "", false
	}
	return data, true
}

// Parse 解析
func Parse(r *http.Request, key string) string {
	if v, ok := r.Form[key]; ok {
		return strings.Join(v, "")
	}
	return ""
}

// Getmd5 获得md5
func Getmd5(name ...string) string {
	var key string
	for i := 0; i < len(name); i++ {
		key += name[i]
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(key))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

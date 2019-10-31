package controller

import (
	"encoding/json"
	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/libs/session/web"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// BackendLogin 登陆
func BackendLogin(w http.ResponseWriter, r *http.Request) {
	data, ok := web.RegularCheck(w, r)
	if !ok {
		return
	}
	s := new(pb.BackendLogin)
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "参数错误"})
		return
	}

	if s.UserName == "ggUser" && s.Password != web.Getmd5("user123456") {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "账号不存在或者密码错误"})
		return
	}

	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	Days := (time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC).Unix()-
		time.Now().Unix())/86400 + 1
	resCookie := http.Cookie{Name: "resCookie", Value: web.Getmd5("resCookie", firstDayOfMonth.Format("2006-01-02")), Path: "/", MaxAge: int(Days) * 86400}
	http.SetCookie(w, &resCookie)
	loginCookie := http.Cookie{Name: "logincookie", Value: web.Getmd5(s.UserName, firstDayOfMonth.Format("2006-01-02")), Path: "/", MaxAge: int(Days) * 86400}
	http.SetCookie(w, &loginCookie)
	// err = authority.GetAllCookies().AddCookie(loginCookie.Value, s.UserName)
	// if err != nil {
	// 	web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "参数错误"})
	// }

	web.HTTPReturnWrite(w, &web.JSONRet{Code: web.StatusOk, Msg: "成功", Data: nil})
}

// UpGetBackendLoginInfo 获取登录信息
func UpGetBackendLoginInfo(w http.ResponseWriter, r *http.Request) {
	data, ok := web.RegularCheck(w, r)
	if !ok {
		return
	}
	s := new(pb.UpGetBackendLoginInfo)
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "参数错误"})
		return
	}
	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)

	loginCookie, err := r.Cookie("logincookie")
	if err != nil {
		logrus.Error("get login cookie err:", err.Error())
	}
	if loginCookie == nil || web.Getmd5("ggUser", firstDayOfMonth.Format("2006-01-02")) != loginCookie.Value {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "账号过期，请重新登录"})
		return
	}

	result := pb.DownGetBackendLoginInfo{
		UserName: "admin",
		Name:     "admin",
	}
	web.HTTPReturnWrite(w, &web.JSONRet{Code: web.StatusOk, Msg: "成功", Data: result})
}

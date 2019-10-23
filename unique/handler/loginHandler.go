package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/libs/mysql"
	"github.com/chaos007/easycome/libs/session/grpc"
	"github.com/chaos007/easycome/libs/session/web"
	"github.com/chaos007/easycome/libs/utils"
	"github.com/chaos007/easycome/model/account"
	"github.com/chaos007/easycome/model/player"
	"github.com/chaos007/easycome/model/verify"

	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

var userMap = &UserMap{
	lock:              new(sync.RWMutex),
	IDToLastLoginTime: map[int64]*account.Account{},
}

//LoginHandler 登录 返回人数最少的agent，并返回服务器列表，数据库可以改装
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data, ok := web.RegularCheck(w, r)
	if !ok {
		return
	}
	s := new(pb.WebUpGetServerInfo)
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "参数错误"})
		return
	}
	sess := mysql.GetEngine().NewSession()
	err = sess.Begin()
	defer sess.Close()
	if err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}

	userInfo := new(account.Account)
	if ok, err := sess.Where("user_name = ? and password = ?", s.UserName, s.Password).Get(userInfo); err != nil {
		logrus.Errorln("LoginHandler get account err:", err)
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	} else if !ok {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "账号密码错误"})
		return
	} else if ok {
		userInfo.Cookie = web.Getmd5(utils.GetRandomString(8))
		userInfo.LastLoginTime = time.Now().Unix()
		if _, err := sess.Where("uid = ?", userInfo.UID).Update(userInfo); err != nil {
			logrus.Errorln("LoginHandler insert err:", err)
			web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
			return
		}
	}

	if err := sess.Commit(); err != nil {
		logrus.Errorln("LoginHandler Commit err:", err)
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}

	result, err := getAllGroupBalanceAgent(userInfo.UID)
	if err != nil || result == nil {
		logrus.Errorln("LoginHandler getAllGroupBalanceAgent err:", err)
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	result.UserID = userInfo.UID
	result.Cookie = userInfo.Cookie
	// result.LastLoginTime = userInfo.LastLoginTime
	web.HTTPReturnWrite(w, &web.JSONRet{Code: 0, Msg: "成功", Data: result})

}

//RegisterHandler 注册
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data, ok := web.RegularCheck(w, r)
	if !ok {
		return
	}
	s := new(pb.WebUpRegister)
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "参数错误"})
		return
	}

	if s.MachineCode == "" {
		for {
			s.MachineCode = utils.GetRandomString(16)
			check1 := &account.Account{}
			if has, err := mysql.GetEngine().Where("machine_code = ?", s.MachineCode).Get(check1); err != nil {
				logrus.Debugln("RegisterHandler err:")
			} else if !has {
				break
			}
		}
	}

	check := &account.Account{}
	if has, err := mysql.GetEngine().Where("machine_code = ?", s.MachineCode).Get(check); err != nil {
		logrus.Debugln("get same machine code err")
	} else if has {
		down, err := getAllGroupBalanceAgent(check.UID)
		if err != nil || down == nil {
			logrus.Errorln("LoginHandler getAllGroupBalanceAgent err:", err)
			web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
			return
		}
		result := &pb.WebDownRegister{
			UserID:      check.UID,
			Password:    check.Password,
			Info:        down.Info,
			Cookie:      check.Cookie,
			MachineCode: s.MachineCode,
		}
		web.HTTPReturnWrite(w, &web.JSONRet{Code: 0, Msg: "成功", Data: result})
		return
	}
	password := utils.GetRandomString(6)
	a := &account.Account{
		Password:    password,
		SignUpTime:  time.Now().Unix(),
		Cookie:      web.Getmd5(password + web.Key),
		UserName:    utils.GetRandomString(8),
		MachineCode: s.MachineCode,
	}
	if _, err := mysql.GetEngine().Insert(a); err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	a.UID = strconv.FormatInt(a.ID+20000, 10)
	if _, err := mysql.GetEngine().Where("id = ?", a.ID).Cols("uid").Update(a); err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	p := &player.Player{
		UID:      a.UID,
		IconID:   1001,
		NickName: a.UserName,
		Gold:     10000,
	}
	if _, err := mysql.GetEngine().Insert(p); err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	down, err := getAllGroupBalanceAgent(a.UID)
	if err != nil || down == nil {
		logrus.Errorln("LoginHandler getAllGroupBalanceAgent err:", err)
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	result := &pb.WebDownRegister{
		UserID:      a.UID,
		Password:    password,
		Info:        down.Info,
		Cookie:      a.Cookie,
		MachineCode: s.MachineCode,
	}
	web.HTTPReturnWrite(w, &web.JSONRet{Code: 0, Msg: "成功", Data: result})
}

//WebUpRegisterPassword 账号注册
func WebUpRegisterPassword(w http.ResponseWriter, r *http.Request) {
	data, ok := web.RegularCheck(w, r)
	if !ok {
		return
	}
	s := new(pb.WebUpRegisterPassword)
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "参数错误"})
		return
	}
	if !utils.IsLetterOrDight(s.UserID) || !utils.IsLetterOrDight(s.Password) {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "账号密码含有非法字符"})
		return
	}
	if len([]rune(s.UserID)) < 6 || len([]rune(s.UserID)) > 16 {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "账号应该在6-16位之间"})
		return
	}
	if len([]rune(s.Password)) < 6 || len([]rune(s.Password)) > 16 {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "密码应该在6-16位之间"})
		return
	}
	if !verify.GetVerifyMap().Check(s.UserID, s.Code) {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "验证码错误或者超时"})
		return
	}
	a := &account.Account{}
	if has, err := mysql.GetEngine().Where("user_name = ?", s.UserID).Get(a); err != nil {
		logrus.Debugln("mysql have same user_name")
	} else if has {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "该账号已被注册"})
		return
	}
	a.UserName = s.UserID
	a.Password = s.Password
	a.SignUpTime = time.Now().Unix()
	a.UserName = utils.GetRandomString(8)
	a.Cookie = web.Getmd5(a.Password + web.Key)

	if _, err := mysql.GetEngine().Insert(a); err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	a.UID = strconv.FormatInt(a.ID+20000, 10)
	if _, err := mysql.GetEngine().Where("id = ?", a.ID).Cols("uid").Update(a); err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	p := &player.Player{
		UID:      a.UID,
		IconID:   1001,
		NickName: a.UserName,
		Gold:     800000,
	}
	if _, err := mysql.GetEngine().Insert(p); err != nil {
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	down, err := getAllGroupBalanceAgent(a.UID)
	if err != nil || down == nil {
		logrus.Errorln("LoginHandler getAllGroupBalanceAgent err:", err)
		web.HTTPReturnWrite(w, &web.JSONRet{Code: web.ParamError, Msg: "系统错误"})
		return
	}
	result := &pb.WebDownRegisterPassword{
		UserID: a.UID,
		Info:   down.Info,
		Cookie: a.Cookie,
	}
	web.HTTPReturnWrite(w, &web.JSONRet{Code: 0, Msg: "成功", Data: result})
}

// AgentCheckUserToLogin 查看玩家是否有登陆信息
func AgentCheckUserToLogin(sess *grpc.Session, data proto.Message) (proto.Message, error) {
	source, ok := data.(*pb.AgentCheckUserToLogin)
	if !ok {
		return nil, nil
	}
	// a := userMap.GetPerson(source.UserID)
	a := &account.Account{}
	result := &pb.LoginCheckPlayerToAgent{
		UserID: source.UserID,
	}

	if has, err := mysql.GetEngine().Where("uid = ? ", source.UserID).Get(a); err != nil {
		logrus.Errorln("LoginHandler get account err:", err)
		return nil, nil
	} else if !has || a.Cookie != source.Cookie {
		result.IsLegal = false
	} else {
		result.IsLegal = true
	}

	return result, nil
}

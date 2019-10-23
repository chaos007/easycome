package grpc

import "sync"

// UserMap rpcstream列表
type UserMap struct {
	list map[int64]*Session
	lock *sync.RWMutex
}

var userMap = &UserMap{
	list: map[int64]*Session{},
	lock: &sync.RWMutex{},
}

// GetUserSession 获得用户的会话
func GetUserSession(userid int64) *Session {
	userMap.lock.RLock()
	defer userMap.lock.RUnlock()
	return userMap.list[userid]
}

// SetUserSession 设置用户的会话
func SetUserSession(s *Session) {
	userMap.lock.Lock()
	defer userMap.lock.Unlock()
	userMap.list[s.UserID] = s
}

// DelUserSession 删除
func DelUserSession(userid int64) {
	userMap.lock.Lock()
	defer userMap.lock.Unlock()
	delete(userMap.list, userid)
}

var serverMap = &ServerMap{
	listSessionKey: map[string]*Session{},
	lock:           &sync.RWMutex{},
}

// ServerMap 服务器rpcstream列表
type ServerMap struct {
	listSessionKey map[string]*Session
	lock           *sync.RWMutex
}

// GetServerSession 获得服务器的会话
func GetServerSession(key string) *Session {
	serverMap.lock.RLock()
	defer serverMap.lock.RUnlock()
	return serverMap.listSessionKey[key]
}

// setUserSession 设置服务器的会话
func setServerSession(s *Session) {
	serverMap.lock.Lock()
	defer serverMap.lock.Unlock()
	serverMap.listSessionKey[s.serverKey] = s
}

// delSession 删除
func delSession(key string) {
	serverMap.lock.Lock()
	defer serverMap.lock.Unlock()
	delete(serverMap.listSessionKey, key)
}

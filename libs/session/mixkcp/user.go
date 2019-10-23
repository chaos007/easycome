package mixkcp

import "sync"

// UserMap rpcstream列表
type UserMap struct {
	listUserID    map[int64]int64
	listSessionID map[int64]*Session
	lock          *sync.RWMutex
}

var userMap = &UserMap{
	listUserID:    map[int64]int64{},
	listSessionID: map[int64]*Session{},
	lock:          &sync.RWMutex{},
}

// GetUserSession 获得用户的会话
func GetUserSession(userid int64) *Session {
	userMap.lock.RLock()
	defer userMap.lock.RUnlock()
	if v, ok := userMap.listUserID[userid]; ok {
		return userMap.listSessionID[v]
	}
	return nil
}

// setUserSession 设置用户的会话
func setUserSession(s *Session) {
	userMap.lock.Lock()
	defer userMap.lock.Unlock()
	userMap.listUserID[s.UserID] = s.SessionID
	userMap.listSessionID[s.SessionID] = s
}

// delSession 删除
func delSession(sessid int64) {
	userMap.lock.Lock()
	defer userMap.lock.Unlock()
	delete(userMap.listSessionID, sessid)
}

// CloseAllSession 关闭所有session
func CloseAllSession() {
	for _, item := range userMap.listSessionID {
		item.SessionClose()
	}
}

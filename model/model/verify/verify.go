package verify

import (
	"sync"
	"time"
)

var l = &ListVerifyNumber{
	lock: new(sync.RWMutex),
	M:    map[string]*PhoneWithVerify{},
}

// GetVerifyMap 获得验证表
func GetVerifyMap() *ListVerifyNumber {
	return l
}

// ListVerifyNumber 验证码表
type ListVerifyNumber struct {
	lock *sync.RWMutex
	M    map[string]*PhoneWithVerify
}

// Add 添加验证码表
func (l *ListVerifyNumber) Add(phone, number string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.clean()
	n, ok := l.M[phone]
	if !ok {
		p := &PhoneWithVerify{
			Phone:  phone,
			Number: number,
			Unix:   time.Now().Add(5 * time.Minute).Unix(),
		}
		l.M[phone] = p
		return
	}
	n.Number = number
	n.Unix = time.Now().Add(5 * time.Minute).Unix()
}

// Check 检查验证码表
func (l *ListVerifyNumber) Check(phone, number string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.clean()
	if n, ok := l.M[phone]; !ok {
		return false
	} else if number != n.Number {
		return false
	} else if time.Now().Unix() > n.Unix {
		return false
	}
	return true
}

func (l *ListVerifyNumber) clean() {
	for _, item := range l.M {
		if time.Now().Unix() > item.Unix {
			delete(l.M, item.Phone)
		}
	}
}

// PhoneWithVerify 手机号带验证表
type PhoneWithVerify struct {
	Phone  string
	Number string
	Unix   int64
}

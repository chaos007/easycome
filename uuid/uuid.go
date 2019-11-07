package uuid

import (
	"strconv"
	"time"
)

var uuidCount = time.Now().UnixNano() / 100

var index int64

// getStringEndNumber GetStringEndNumber
func getStringEndNumber(s string) int64 {
	index := len(s) - 1
	for i := len(s) - 1; i >= 0; i-- {
		if uint8(s[i]) < 48 || uint8(s[i]) > 57 {
			break
		} else {
			index = i
		}
	}

	i, err := strconv.ParseInt(s[index:], 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// InitIndex 一组服务器对应数据库可以标识自己的id，最多有100个
func InitIndex(s string) bool {
	index = getStringEndNumber(s)
	if index >= 100 || index < 0 {
		return false
	}
	return true
}

// GenUUID 生成uuid
func GenUUID() int64 {
	now := time.Now().UnixNano() / 100
	if now <= uuidCount {
		uuidCount++
		return uuidCount + index
	}
	uuidCount = time.Now().UnixNano() / 100
	return uuidCount + index
}

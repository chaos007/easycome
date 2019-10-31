package utils

import (
	"math/rand"
	"time"
	"unicode"
)

// GetRandomString 获得随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < l; i++ {
		r := int(time.Now().Unix()) * rand.Int()
		if r < 0 {
			r *= -1
		}
		index := r % (len(bytes))
		result = append(result, bytes[index])
	}
	return string(result)
}

// IsLetterOrDight 是否由字符后者数字组成
func IsLetterOrDight(str string) bool {
	for _, r := range str {
		if !(unicode.IsDigit(r) || unicode.IsLetter(r)) {
			return false
		}
	}
	return true
}

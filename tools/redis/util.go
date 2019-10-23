package redis

import (
	"strconv"
)

func getBool(result string) bool {
	if result == "OK" {
		return true
	} else {
		return false
	}
}

func float64ToStr(floatNumber float64) string {
	return strconv.FormatFloat(floatNumber, 'f', -1, 64)
}

func strToFloat64(floatString string) float64 {
	res, _ := strconv.ParseFloat(floatString, 64)
	return res
}

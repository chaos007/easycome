package service

import (
	"github.com/chaos007/easycome/libs/config"
)

// Info the detail of service
type Info struct {
	IP string
}

var me *Service

// GetMe 获得自己
func GetMe() *Service {
	return me
}

// Service service
type Service struct {
	Name        string
	IP          string
	stop        chan error
	General     config.General
	servicePath string
	infoPlace   string
	serviceType string
	group       string
	version     string
}

package handler

import "github.com/chaos007/easycome/libs/session/grpc"

// ProtocolInit Init
func ProtocolInit() {
	grpc.RegisterProtocol("pb.AgentCheckUserToLogin", AgentCheckUserToLogin)

}

package handle

import (
	"github.com/chaos007/easycome/agent/handle/controller"
	"github.com/chaos007/easycome/libs/session/grpc"
	"github.com/chaos007/easycome/libs/session/mixkcp"
)

// Init Init
func Init() {
	grpc.RegisterProtocol("pb.TestToAgent", controller.TestToAgent)
	mixkcp.RegisterProtocol("pb.TestStruct", controller.TestStruct)
	mixkcp.RegisterProtocol("pb.UpPing", controller.UpPing)

	mixkcp.RegisterProtocol("pb.PlayerLogin", controller.PlayerLogin)
	mixkcp.RegisterProtocol("pb.LoginCheckPlayerToAgent", controller.LoginCheckPlayerToAgent)
	grpc.RegisterProtocol("pb.LoginKickPlayerToAgent", controller.LoginKickPlayerToAgent)

}

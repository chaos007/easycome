package handle

import (
	"net/http"

	"github.com/chaos007/easycome/center/handle/controller"
	"github.com/chaos007/easycome/libs/session/grpc"
)

// Init Init
func Init() {
	grpc.RegisterProtocol("pb.UpGetConfigToCenter", controller.UpGetConfigToCenter)

}

// HTTPInit 后台
func HTTPInit() {
	http.HandleFunc("/BackendLogin", controller.BackendLogin)
	http.HandleFunc("/UpGetBackendLoginInfo", controller.UpGetBackendLoginInfo)
}

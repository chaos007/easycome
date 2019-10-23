package controller

import (
	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/libs/session/grpc"
	"github.com/chaos007/easycome/model/configdata"
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// UpGetConfigToCenter test
func UpGetConfigToCenter(sess *grpc.Session, data proto.Message) (proto.Message, error) {
	_, ok := data.(*pb.UpGetConfigToCenter)
	if !ok {
		return nil, nil
	}

	design := configdata.GetDesign()

	configString, err := json.Marshal(design)
	if err != nil {
		logrus.Errorln("json marshal configdata err:", err)
	}

	result := &pb.SyncConfigDataToAll{
		Data: string(configString),
	}

	return result, nil
}

package controller

import (
	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/enum"
	"github.com/chaos007/easycome/libs/session/grpc"
	"github.com/chaos007/easycome/model/configdata"
	"github.com/chaos007/easycome/model/player"

	"github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// TestToGame test
func TestToGame(sess *grpc.Session, data proto.Message) (proto.Message, error) {
	source := data.(*pb.TestToGame)

	logrus.Debugf("game receve data:%#v", source.FieldA)

	// result := &pb.TestStruct{
	// 	FieldA: "44444",
	// }

	to := &pb.TestToAgent{
		FieldA: "2222",
	}

	err := sess.SendToSteamServer(enum.ServerTypeAgent, to)
	if err != nil {
		logrus.Debugln("-----------err:", err)
	}

	return nil, nil
}

// SyncConfigDataToAll test
func SyncConfigDataToAll(sess *grpc.Session, data proto.Message) (proto.Message, error) {
	source, ok := data.(*pb.SyncConfigDataToAll)

	if !ok {
		return nil, nil
	}

	err := configdata.UpdateDesignConfig(source.Data)
	if err != nil {
		logrus.Errorln("SyncConfigDataToAll update design config err:", err)
	}
	return nil, nil
}

// PlayerLoginToGame test
func PlayerLoginToGame(sess *grpc.Session, data proto.Message) (proto.Message, error) {
	source, ok := data.(*pb.PlayerLoginToGame)

	if !ok {
		return nil, nil
	}

	p, err := player.NewPlayer(source.UserID)

	if err != nil {
		logrus.Errorln("PlayerLoginToGame get player err:", err)
	}

	sess.Player = p

	result := &pb.DownPlayerLogin{
		UserID: p.UID,
		// IsSucceed:   true,
		// Gold:        p.Gold,
		// Money:       p.Money,
		// Level:       p.Level,
		// Exp:         p.Exp,
		// Name:        p.Name,
		// Power:       p.Power,
		// Stamina:     p.Stamina,
		// Agility:     p.Agility,
		// Attack:      p.Attack,
		// Defence:     p.Defence,
		// HP:          p.HP,
		// Dodge:       p.Dodge,
		// AttackSpeed: p.AttackSpeed,
		// AttackRate:  p.AttackRate,
		// HitRate:     p.HitRate,
		// RoleID:      p.RoleID,
	}
	return result, nil
}

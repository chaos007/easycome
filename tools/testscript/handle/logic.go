package handle

import (
	"fmt"
	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/tools/testscript/player"

	"github.com/golang/protobuf/proto"
)

// DownPing DownPing
func DownPing(content proto.Message, p *player.Player) (ret proto.Message) {
	msg := content.(*pb.DownPing)
	fmt.Println("------DownPing:", msg)
	//p.Done("DownPing")
	return
}

// DownCreateBattleRoom DownCreateBattleRoom
func DownCreateBattleRoom(content proto.Message, p *player.Player) (ret proto.Message) {
	msg := content.(*pb.DownCreateBattleRoom)
	fmt.Println("------DownCreateBattleRoom:", msg)
	p.Done(msg.RoomID)
	return
}

// DownJoinBattleRoom DownJoinBattleRoom
func DownJoinBattleRoom(content proto.Message, p *player.Player) (ret proto.Message) {
	msg := content.(*pb.DownJoinBattleRoom)
	// fmt.Println("------DownJoinBattleRoom:", msg)
	if msg != nil || msg.PlayerList != nil {
		for _, item := range msg.PlayerList {
			fmt.Println("------DownJoinBattleRoom:", item.UserID)
		}
	}
	p.Done("DownJoinBattleRoom")
	return
}

// DownJoinBattleRoomToOtherPlayer DownJoinBattleRoomToOtherPlayer
func DownJoinBattleRoomToOtherPlayer(content proto.Message, p *player.Player) (ret proto.Message) {
	msg := content.(*pb.DownJoinBattleRoomToOtherPlayer)
	fmt.Println("------DownJoinBattleRoomToOtherPlayer:", msg)
	// p.Done("DownJoinBattleRoomToOtherPlayer")
	return
}

// DownStartBattle DownStartBattle
func DownStartBattle(content proto.Message, p *player.Player) (ret proto.Message) {
	msg := content.(*pb.DownStartBattle)
	fmt.Println("------DownStartBattle:", msg)
	// p.Done("DownStartBattle")
	return
}

// DownStartBattleToOtherPlayer DownStartBattleToOtherPlayer
func DownStartBattleToOtherPlayer(content proto.Message, p *player.Player) (ret proto.Message) {
	msg := content.(*pb.DownStartBattleToOtherPlayer)
	fmt.Println("------DownStartBattleToOtherPlayer:", msg)
	// p.Done("DownStartBattleToOtherPlayer")
	return
}

// DownGameStartToAll DownGameStartToAll
func DownGameStartToAll(content proto.Message, p *player.Player) (ret proto.Message) {
	msg := content.(*pb.DownGameStartToAll)
	fmt.Println("------DownGameStartToAll:", msg)
	p.Done("DownGameStartToAll")
	return
}

// DownPlayerAction DownPlayerAction
func DownPlayerAction(content proto.Message, p *player.Player) (ret proto.Message) {
	// msg := content.(*pb.DownPlayerAction)
	// fmt.Println("------DownPlayerAction:", msg)
	// p.Done("DownPlayerAction")
	return
}

// GameFrame GameFrame
func GameFrame(content proto.Message, p *player.Player) (ret proto.Message) {
	// msg := content.(*pb.GameFrame)
	// fmt.Println("------GameFrame:", msg)
	// p.Done("GameFrame")
	return
}

// EmptyGameFrame EmptyGameFrame
func EmptyGameFrame(content proto.Message, p *player.Player) (ret proto.Message) {
	// msg := content.(*pb.EmptyGameFrame)
	// fmt.Println("------EmptyGameFrame:", msg)
	// p.Done("EmptyGameFrame")
	return
}

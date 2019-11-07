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

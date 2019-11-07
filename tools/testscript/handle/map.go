package handle

import (
	"fmt"

	"github.com/chaos007/easycome/libs/msgmeta"
	"github.com/chaos007/easycome/tools/testscript/player"

	"github.com/golang/protobuf/proto"
)

// Handler 处理数据
func Handler(msgID int32, data []byte, p *player.Player) (ret proto.Message) {
	v, ok := msgmeta.MessageMetaByID(msgID)
	if !ok {
		fmt.Printf("msgID:%d  not found!", msgID)
		return nil
	}
	callback, ok := handler[v]
	if !ok {
		fmt.Printf("proto callback name:%s not found!", v)
		return nil
	}
	pdata, err := msgmeta.ParsePacket(msgID, data)
	if err != nil {
		fmt.Printf("ParsePacket,msgID:%d", msgID)
		return nil
	}
	ret = callback(pdata, p)
	return nil
}

//大厅网关协议
var handler = map[string]func(proto.Message, *player.Player) proto.Message{
	"pb.DownPing":         DownPing,
	"pb.DownPlayerAction": DownPlayerAction,
	"pb.EmptyGameFrame":   EmptyGameFrame,
	"pb.GameFrame":        GameFrame,
}

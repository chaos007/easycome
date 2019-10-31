package client

import (
	"github.com/chaos007/easycome/tools/testscript/player"
)

// SendHandler SendHandler
var SendHandler = map[int32]func(*Client, *player.Player) interface{}{}

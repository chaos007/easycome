package player

import (
	"github.com/chaos007/easycome/libs/mysql"
	"github.com/chaos007/easycome/model/configdata"
)

// Player 玩家
type Player struct {
	ID        int64  `xorm:"autoincr id"`
	UID       string `xorm:"index uid"`
	IsMale    bool
	NickName  string
	Signature string //签名
	Gold      int64
	Exp       int64
	VipLevel  int32
	Level     int32
}

// NewPlayer NewPlayer
func NewPlayer(playerID string) (*Player, error) {
	p := new(Player)
	if ok, err := mysql.GetEngine().Where("uid = ?", playerID).Get(p); err != nil {
		return nil, err
	} else if !ok {
		p.UID = playerID
		if _, err := mysql.GetEngine().Insert(p); err != nil {
			return nil, err
		}
	}
	return p, p.Init()
}

// Init 初始化
func (p *Player) Init() (err error) {
	return nil
}

// AddExp 增加经验
func (p *Player) AddExp(exp int64) {
	level := configdata.GetDesign().LevelData
	p.Exp += exp
	for {
		var isChange bool
		for _, l := range level {
			if l.Level == int64(p.Level) {
				if p.Exp > l.Exp {
					p.Exp -= l.Exp
					p.Level++
					isChange = true
				}
			}
		}
		if !isChange {
			break
		}
	}
}

//AddGold 增加或者减少金币
func (p *Player) AddGold(gold int64) bool {
	if p.Gold+gold < 0 {
		return false
	}
	p.Gold += gold
	return true
}

// Save 保存玩家数据
func (p *Player) Save() (err error) {
	return nil
}

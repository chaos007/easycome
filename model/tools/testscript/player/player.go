package player

import (
	"fmt"
	"sync"
)
import (
	"github.com/chaos007/easycome/tools/testscript/types"
)

// PlayersManager Players map[int64]*Player
var PlayersManager = NewPlayers()

// Players player
type Players struct {
	data map[int64]*Player
	lock sync.RWMutex
}

// NewPlayers NewPlayers
func NewPlayers() *Players {
	return &Players{
		data: make(map[int64]*Player),
	}
}

// Get Get
func (p *Players) Get(id int64) *Player {
	p.lock.Lock()
	defer p.lock.Unlock()
	if a, ok := p.data[id]; ok {
		return a
	}
	return nil
}

// Add Add
func (p *Players) Add(a *Player) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data[a.ID] = a
}

// Delete Delete
func (p *Players) Delete(id int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.data, id)
}

// Player Player
type Player struct {
	ID               int64
	NickName         string
	Level            int32
	CurrentExp       int32
	Gold             int32
	CurrentGameState int32
	SelectHeroID     int32
	Session          *types.Session
}

// SetSession SetSession
func (p *Player) SetSession(sess *types.Session) {
	p.Session = sess
}

// Done Done
func (p *Player) Done(method interface{}) {
	fmt.Printf("%s done\n", method)
	p.Session.Step <- method
}

// Wait Wait
func (p *Player) Wait() interface{} {
	method := <-p.Session.Step
	fmt.Printf("wait for after %s\n", method)
	return method
}

// GetPlayer GetPlayer
func GetPlayer(userID int64) *Player {
	if p := PlayersManager.Get(userID); p != nil {
		return p
	}
	p := NewPlayer()
	p.ID = userID
	PlayersManager.Add(p)
	return p
}

// NewPlayer NewPlayer
func NewPlayer() *Player {
	return &Player{
		NickName: "",
		// SaveList: make([]pb.PBSaveListItem, 0),
	}
}

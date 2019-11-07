package interfacer

// Player 玩家接口
type Player interface {
	Save()
	NewPlayer(string) error
}

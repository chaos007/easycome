package client

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/tools/testscript/player"
	"github.com/chaos007/easycome/tools/testscript/socket"
	"github.com/chaos007/easycome/tools/testscript/types"
)

// Client 客户端
type Client struct {
	Session *types.Session
	Socket  *socket.Socket
}

// LogFile LogFile
var LogFile *os.File

// SetAgentHost 设置头
func (c *Client) SetAgentHost(host string) {
	c.Socket.SetHost(host)
}

// Do 通信
func (c *Client) Do(die chan bool, p *player.Player) (err error) {
	if err = c.Socket.Dail(); err != nil {
		fmt.Println("--------err:", err)
		fmt.Printf("-------Dail server failed!")
	} else {
		go c.Socket.RecvThread(c.Session, die, p)
	}
	return
}

// Play Play
func (c *Client) Play(p *player.Player, battleType int32) string {
	return "Down"
}

func uppingtimer(c *Client) {
	timer2 := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-timer2.C:
			c.Socket.Send(&pb.UpPing{})
		}
	}
}

func gameoptimer(c *Client, p *player.Player) {
	timer1 := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer1.C:
			// gameopsend(c)
		case msg := <-p.Session.Step:
			if msg == "DownGameOver" {
				// break
				return
			}
		}
	}
}

// // Getgameop 模拟游戏操作
// func Getgameop() []byte {
// 	data := &syncInfosPB.ControlData{}
// 	vda := Getrandbool()
// 	vdd := Getrandbool()
// 	vds := Getrandbool()
// 	vdw := Getrandbool()
// 	vdf := Getrandbool()
// 	vs1 := Getrandbool()
// 	vs2 := Getrandbool()
// 	vs3 := Getrandbool()
// 	vmx := Getrandint()
// 	vmy := Getrandint()
// 	vmd := Getrandbool()
// 	vbeid := Getrandint()
// 	vseid := Getrandint()
// 	vlk := Getrandint()
// 	vdr := Getrandbool()
// 	vdv := Getrandbool()
// 	//vsurrender := Getrandint()
// 	vminmapx := Getrandint()
// 	vminmapy := Getrandint()
// 	data.Da = &vda
// 	data.Dd = &vdd
// 	data.Ds = &vds
// 	data.Dw = &vdw
// 	data.Df = &vdf
// 	data.S1 = &vs1
// 	data.S2 = &vs2
// 	data.S3 = &vs3
// 	data.Mx = &vmx
// 	data.My = &vmy
// 	data.Md = &vmd
// 	data.BuyEquipmentId = &vbeid
// 	data.SellEquipmentId = &vseid
// 	data.LearnSkillId = &vlk
// 	data.Chat = nil
// 	data.Dr = &vdr
// 	data.Dv = &vdv
// 	//data.Surrender = &vsurrender
// 	data.MinMapx = &vminmapx
// 	data.MinMapy = &vminmapy
// 	r, err := proto.Marshal(data)
// 	if err != nil {
// 		log.Fatal("marshaling error: ", err)
// 	}
// 	return r
// }

// Getrandbool Getrandbool
func Getrandbool() bool {
	r := rand.Intn(2)
	if r == 0 {
		return true
	}
	return false
}

// Getrandint Getrandint
func Getrandint() int32 {
	return rand.Int31n(10)
}

// TestRandOrder 模拟一系列指令操作
func TestRandOrder(addr string, die chan bool, id int) {
	fmt.Printf("playerID:%d", id)
	Client := NewClient()
	p := player.GetPlayer(int64(id))
	p.ID = int64(id)
	Client.Session.UserID = p.ID
	p.SetSession(Client.Session)
	var err error

	// fileName := "send.txt"
	// os.Remove(fileName)
	// LogFile, err := os.Create(fileName)
	// defer LogFile.Close()
	// if err != nil {
	// 	LogFile.WriteString("open file error !")
	// }

	Client.Socket.SetHost(addr)
	if err = Client.Do(die, p); err != nil {
		println("err:", err)
		return
	}
	go p.Wait()

	go blackboxtimer(Client, p)
}

func blackboxtimer(c *Client, p *player.Player) {
	timer3 := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer3.C:
			orderid := rand.Int31n(13) + 18
			callback, ok := SendHandler[orderid]
			if !ok {
				println("order err:", orderid)
				return
			}
			callback(c, p)
		}
	}
}

// NewClient NewClient
func NewClient() *Client {
	return &Client{
		Session: types.NewSession(),
		Socket:  socket.NewSocket(),
	}
}

// UpPlayerAction UpPlayerAction
func (c *Client) UpPlayerAction() {
	i := 0
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			if i%100 == 0 {
				fmt.Printf("-------------rece from tick:%d,userid:%d\n", i, c.Session.UserID)
			}
			i++
			c.Socket.Send(&pb.UpPlayerAction{
				RoomID:  1,
				Message: []byte(fmt.Sprintf("this is userid :%d,tick:%d", c.Session.UserID, i)),
				UserID:  c.Session.UserID,
			})
		}
	}

}

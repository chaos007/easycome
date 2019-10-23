package main

import (
	"fmt"
	"github.com/chaos007/easycome/data/pb"
	"github.com/chaos007/easycome/tools/testscript/client"
	"github.com/chaos007/easycome/tools/testscript/player"
	"sync/atomic"
	"testing"
	"time"
)

var testHost = "218.107.220.124:11000"

func TestServer(t *testing.T) {
	ch := make(chan int)
	for i := 0; i < 2000; i += 100 {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
		go func(i int) {
			room(int64(i))
		}(i)
	}
	<-ch
}

func room(i int64) {
	ch := make(chan int)

	var dieClient = make(chan bool)
	c := client.NewClient()
	c.Session.UserID = int64(i)
	p := player.GetPlayer(int64(i))
	p.ID = int64(i)
	p.SetSession(c.Session)
	c.SetAgentHost(testHost)
	c.Do(dieClient, p)

	c.Socket.Send(&pb.UpCreateBattleRoom{
		Player: &pb.BattlePlayer{
			UserID: int64(i),
		},
	})
	roomID := p.Wait()
	fmt.Println("----roomID", roomID)
	list := otherPlayerJoinRoom(roomID.(int64), i)
	time.Sleep(5 * time.Second) //等待玩家加入房间
	fmt.Println("----start")
	c.Socket.Send(&pb.UpStartBattle{
		BattleRoomID: roomID.(int64),
		UserID:       int64(i),
	})
	go func() {
		msg := p.Wait()
		if msg == "DownGameStartToAll" {
			go c.UpPlayerAction()
		}
	}()

	for _, item := range list {
		go func(c *client.Client) {
			c.Socket.Send(&pb.UpStartBattle{
				BattleRoomID: roomID.(int64),
				UserID:       c.Session.UserID,
			})
			o := player.GetPlayer(c.Session.UserID)
			msg := o.Wait()
			if msg == "DownGameStartToAll" {
				go c.UpPlayerAction()
			}
		}(item)
	}

	<-ch
}

func otherPlayerJoinRoom(id, idbegin int64) []*client.Client {
	result := []*client.Client{}
	for i := idbegin + 1; i < idbegin+10; i++ {
		var dieClient = make(chan bool)
		c := client.NewClient()
		c.Session.UserID = int64(i)
		p := player.GetPlayer(int64(i))
		p.ID = int64(i)
		p.SetSession(c.Session)
		c.SetAgentHost(testHost)
		c.Do(dieClient, p)
		c.Socket.Send(&pb.UpJoinBattleRoom{
			BattleRoomID: id,
			Player: &pb.BattlePlayer{
				UserID: int64(i),
			},
		})
		result = append(result, c)
		p.Wait()
	}
	return result
}

func TestAtomic(t *testing.T) {

	var i int32 = -1 * (1<<30 + 555)
	var addr = &i
	fmt.Println(atomic.AddInt32(addr, -1<<30))
	result := i + 1<<30
	fmt.Println(result)
}

func TestConstMap(t *testing.T) {
	// const mapList = map[int]int{}
	// mapList[1] = 1
}

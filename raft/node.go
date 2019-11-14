package raft

import "time"

// 三种状态
const (
	StateLeader       = 1
	StateFollower     = 2
	StateCandidate    = 3
	StatePreCandidate = 4
)

// Node 当前自己的节点
type Node struct {
	State        int
	LastTermUnix int64
	TermMill     int64
	TeamIndex    uint64
}

// RecvHeartBeat 收到心跳
func (n *Node) RecvHeartBeat(term uint64) {
	if n.State == StateLeader { //自己作为leader收到心跳包
		if term < n.TeamIndex { //小于自己的选举期
			return
		}
		if term >= n.TeamIndex { //大于自己的选举期，自己要变成follower

		}
	}
	n.State = StateFollower //其他两个状态变为follower
	n.LastTermUnix = time.Now().UnixNano()
}

// Colck 时间周期
func (n *Node) Colck() {
	t := time.NewTicker(10 * time.Millisecond)

	go func() {
		select {
		case <-t.C:
			if n.State == StateLeader { //状态是领导者
				break
			}
			if (time.Now().UnixNano() - n.LastTermUnix) < n.TermMill*1000*1000 { //未到选举期
				break
			}
		}
	}()

}

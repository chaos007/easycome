package uuid

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var m = map[int64]bool{}

var lock = sync.RWMutex{}

func TestString(t *testing.T) {
	fmt.Println(getStringEndNumber("game123"))
}

func TestUUID(t *testing.T) {

	begin := time.Now().UnixNano()

	go func() {
		v := GenUUID()
		lock.Lock()
		if _, ok := m[v]; ok {
			fmt.Println("---------")
		}
		m[v] = true
		lock.Unlock()
	}()

	go func() {
		v := GenUUID()
		lock.Lock()
		if _, ok := m[v]; ok {
			fmt.Println("---------")
		}
		m[v] = true
		lock.Unlock()
	}()
	for i := 0; i < 1000; i++ {
		v := GenUUID()
		lock.Lock()
		if _, ok := m[v]; ok {
			fmt.Println("---------")
		}
		m[v] = true
		lock.Unlock()
	}

	end := time.Now().UnixNano()

	fmt.Println("last :", end-begin)

	// loginch = make(chan int)

	// go httpTest()

	// <-loginch

	// tcpTest(t)

}

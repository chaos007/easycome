package redis

import (
	rg "github.com/garyburd/redigo/redis"
	"time"
)

type (
	PoolConnection interface {
		RedisCommands
		PoolCommands
	}

	ConnectionPoolConfig struct {
		MaxIdle     int
		MaxActive   int
		IdleTimeout time.Duration
	}
)

var DefaultConnectionPoolConfig = ConnectionPoolConfig{
	MaxIdle:     60,
	MaxActive:   100,
	IdleTimeout: 30,
}

func CreatePool(serverAddr, auth, db string, poolConfig ConnectionPoolConfig) (PoolConnection, error) {
	dialer := func() (rg.Conn, error) {
		return Connect(serverAddr, auth, db)
	}

	tester := func(c rg.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	}

	config := getConnectionPoolConfig(&poolConfig)

	pool := &rg.Pool{
		MaxIdle:      config.MaxIdle,
		IdleTimeout:  config.IdleTimeout,
		MaxActive:    config.MaxActive,
		Dial:         dialer,
		TestOnBorrow: tester,
	}

	con := &connection{p: pool}

	return con, nil
}

func CreatePoolUri(uri string, poolConfig ConnectionPoolConfig) (PoolConnection, error) {
	dialer := func() (rg.Conn, error) {
		return ConnectUrl(uri)
	}

	tester := func(c rg.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	}

	config := getConnectionPoolConfig(&poolConfig)

	pool := &rg.Pool{
		MaxIdle:      config.MaxIdle,
		IdleTimeout:  config.IdleTimeout,
		Dial:         dialer,
		TestOnBorrow: tester,
	}

	con := &connection{p: pool}

	return con, nil
}

func getConnectionPoolConfig(config *ConnectionPoolConfig) ConnectionPoolConfig {
	if config.IdleTimeout == 0 || config.MaxIdle == 0 || config.MaxActive == 0 {
		return DefaultConnectionPoolConfig
	} else {
		return *config
	}
}

func (con *connection) GetConnection() (PoolConnection, error) {
	c := con.p.Get()
	resCon := &connection{c: c}
	return resCon, nil
}

func (con *connection) ActiveCount() int {
	return con.p.ActiveCount()
}

func (con *connection) Release() {
	con.c.Close()
}

func (con *connection) PoolClose() {
	con.p.Close()
}

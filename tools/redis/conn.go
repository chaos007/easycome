package redis

import (
	"fmt"
)

type Connection interface {
	RedisCommands
	ConnectionCommands
}

func CreateConnection(serverAddr, auth, db string) (Connection, error) {
	c, e := Connect(serverAddr, auth, db)

	if e != nil {
		fmt.Println(e)

		return nil, e
	}

	con := &connection{c: c}

	return con, nil
}

func CreateConnectionUri(uri string) (Connection, error) {
	c, e := ConnectUrl(uri)

	if e != nil {
		fmt.Println(e)

		return nil, e
	}

	con := &connection{c: c}

	return con, nil
}

func (con *connection) Close() bool {
	err := con.c.Close()

	if err != nil {
		return false
	}

	return true
}

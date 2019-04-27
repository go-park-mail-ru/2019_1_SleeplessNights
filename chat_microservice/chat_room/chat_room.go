package chat_room

import "golang.org/x/net/websocket"

var chat *chatRoom

const (
	maxConnections = 100
)

type chatRoom struct {
	maxConnections int64
	ConnectionPool []*websocket.Conn
}

func init() {
	chat = &chatRoom{
		maxConnections: maxConnections,
	}
}

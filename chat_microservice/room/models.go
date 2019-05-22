package room

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Talker struct {
	Conn       *websocket.Conn
	Nickname   string
	AvatarPath string
	Id         uint64
}

type Message struct {
	Title   string  `json:"title"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Text  string `json:"text,omitempty"`
	Since uint64 `json:"since,omitempty"`
}

type ResponseMessage struct {
	Nickname   string `json:"nickname"`
	AvatarPath string `json:"avatarPath"`
	Id         uint64 `json:"id"`
	Text       string `json:"text"`
}

type room struct {
	maxConnections int64
	Id             uint64
	usersPool      map[uint64]*Talker
	mx             sync.Mutex
}

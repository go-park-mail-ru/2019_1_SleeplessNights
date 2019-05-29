package room_manager

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

type message struct {
	Title   string  `json:"title"`
	Payload payload `json:"payload"`
}

type payload struct {
	Text  string `json:"text,omitempty"`
	Since uint64 `json:"since,omitempty"`
}
//easyjson:json
type responseMessage struct {
	Nickname   string `json:"nickname"`
	AvatarPath string `json:"avatarPath"`
	Id         uint64 `json:"id"`
	Text       string `json:"text"`
}

type room struct {
	MaxConnections uint64
	Id             uint64
	TalkersPool    map[uint64]*Talker
	AccessArray    []uint64
	mx             sync.Mutex
}

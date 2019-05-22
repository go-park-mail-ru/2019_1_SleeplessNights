package database

type talker struct {
	Id         uint64 `sql:"id"`
	Uid        uint64 `sql:"uid"`
	NickName   string `sql:"nickname"`
	AvatarPath string `sql:"avatar_path"`
}

type message struct {
	Id       uint64 `sql:"id"`
	Payload  string `sql:"payload"`
	TalkerId uint64 `sql:"author_id"`
	RoomId   uint64 `sql:"room_id"`
}

type room struct {
	Id      uint64   `sql:"id"`
	Talkers []uint64 `sql:"authors"`
}

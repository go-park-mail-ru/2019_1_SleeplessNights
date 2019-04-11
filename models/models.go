package models

import (
	"time"
)

type User struct {
	ID         uint          `json:"-"`
	Email      string        `json:"email"`
	Password   []byte        `json:"-"`
	Salt       []byte        `json:"-"`
	Won        uint          `json:"won"`
	Lost       uint          `json:"lost"`
	PlayTime   time.Duration `json:"play_time"`
	Nickname   string        `json:"nickname"`
	AvatarPath string        `json:"avatar_path"`
}

//type UserPk struct {
//	ID    uint
//	Email string
//}

type Question struct {
	ID      uint
	Answers []string
	Correct int
	Text    string
	PackID  uint
}

type Pack struct {
	Question []Question
	Theme    string
}

var idSource uint
//var Users map[string]User
//var UserKeyPairs map[uint]string
//
//func init() {
//	Users = make(map[string]User, 0)
//	UserKeyPairs = make(map[uint]string, 0)
//}

func MakeID() uint {
	idSource++
	return idSource
}

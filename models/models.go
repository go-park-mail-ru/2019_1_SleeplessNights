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

type Answer struct {
	first  string
	second string
	third  string
	forth  string
}

type Question struct {
	ID      uint
	Answers []string
	Correct int
	Text    string
	PackID  uint
	Theme   string
}

type Pack struct {
	Question []Question
	Theme    string
}

var idSource uint

func MakeID() uint {
	idSource++
	return idSource
}

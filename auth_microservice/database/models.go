package database

import "time"

type User struct {
	ID         uint64        `json:"-"`
	Email      string        `json:"email"`
	Password   []byte        `json:"-"`
	Salt       []byte        `json:"-"`
	Won        uint          `json:"won"`
	Lost       uint          `json:"lost"`
	PlayTime   time.Duration `json:"play_time"`
	Nickname   string        `json:"nickname"`
	AvatarPath string        `json:"avatar_path"`
}
package models

import (
	"log"
	"os"
	"time"
)

type User struct {
	ID        uint
	Email     string
	Password  []byte//Длина хеш-суммы по алгоритму SHA512
	Salt      []byte
	SessionID uint
	ProfileID uint
	BestScore uint
}

type UserPk struct {
	ID    uint
	Email string
}

type Profile struct {
	ID       uint
	Nickname string
	AvatarID uint
}

type Avatar struct {
	ID   uint
	Path string
}

var idSource uint
var Users map[string]User
var UserKeyPairs map[uint]string
var Profiles map[uint]Profile
var Avatars map[uint]Avatar

func init() {
	Users = make(map[string]User, 0)
	UserKeyPairs = make(map[uint]string, 0)
	Profiles = make(map[uint]Profile, 0)
	Avatars = make(map[uint]Avatar, 0)

	err := os.Setenv("SaltLen", "16")//16 байт (128 бит), как в современных UNIX системах
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("SessionLifeLen", time.Hour.String())
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("ServerID", "MyConfidantServer")
	if err != nil {
		log.Fatal(err)
	}
}

func MakeID() uint {
	//TODO make thread-safe
	idSource++
	return idSource
}

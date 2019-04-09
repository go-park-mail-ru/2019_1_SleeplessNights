package database

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
)

var db *dbManager

type dbManager struct {
	users        map[string]models.User
	userKeyPairs map[uint]string
}

func init() {
	db = &dbManager{
		users:        make(map[string]models.User, 0),
		userKeyPairs: make(map[uint]string, 0),
	}
}

func GetUserViaID(userID uint) (models.User, bool) {
	user, found := db.users[db.userKeyPairs[userID]]
	return user, found
}

func GetUserViaEmail(email string) (models.User, bool) {
	user, found := db.users[email]
	return user, found
}

func GetUserKeyPair(userID uint) string {
	email:= db.userKeyPairs[userID]
	return email
}

func AddIntoUsers(user models.User, email string)  {
	db.users[email] = user
}

func AddIntoUserKeyPairs(email string, id uint) {
	db.userKeyPairs[id] = email
}

func DeleteIntoUsers(email string) {
	delete(db.users, email)
}

func DeleteIntoUserKeyPairs(id uint) {
	delete(db.userKeyPairs, id)
}

func GetLenUsers() int {
	return len(db.users)
}

func GetUsers() map[string]models.User {
	duplicate := db.users
	return duplicate
}

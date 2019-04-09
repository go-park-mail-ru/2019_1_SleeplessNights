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

func GetInstance() *dbManager {
	return db
}

func(db *dbManager) GetUserViaID(userID uint) (models.User, bool) {
	user, found := db.users[db.userKeyPairs[userID]]
	return user, found
}

func(db *dbManager) GetUserViaEmail(email string) (models.User, bool) {
	user, found := db.users[email]
	return user, found
}

func(db *dbManager) GetUserKeyPair(userID uint) string {
	email:= db.userKeyPairs[userID]
	return email
}

func(db *dbManager) AddIntoUsers(user models.User, email string)  {
	db.users[email] = user
}

func(db *dbManager) AddIntoUserKeyPairs(email string, id uint) {
	db.userKeyPairs[id] = email
}

func(db *dbManager) DeleteIntoUsers(email string) {
	delete(db.users, email)
}

func(db *dbManager) DeleteIntoUserKeyPairs(id uint) {
	delete(db.userKeyPairs, id)
}

func(db *dbManager) GetLenUsers() int {
	return len(db.users)
}

func(db *dbManager) GetUsers() map[string]models.User {
	duplicate := db.users
	return duplicate
}

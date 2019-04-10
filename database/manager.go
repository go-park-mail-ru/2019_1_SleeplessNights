package database

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
)

//const (
//	host     = "localhost"
//	port     = 5432
//	user     = "postgres"
//	password = "1209qawsed"
//	dbName   = "postgres"
//)

const (
	host     = ""
	port     = 0
	user     = ""
	password = ""
	dbName   = ""
)

var db *dbManager

type dbManager struct {
	dateBase *sql.DB
}

func OpenConnection() (err error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	dateBase, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}

	err = dateBase.Ping()
	if err != nil {
		return
	}

	db = &dbManager{
		dateBase: dateBase,
	}

	return
}

func CloseConnection() (err error) {
	err = db.dateBase.Close()
	return
}

func GetInstance() *dbManager {
	return db
}

func (db *dbManager) GetUserViaID(userID uint) (user models.User, found bool, err error) {

	row := db.dateBase.QueryRow(
		`SELECT * FROM public.users WHERE id = $1`, userID)
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Salt, &user.Won, &user.Lost, &user.PlayTime, &user.Nickname,
		&user.AvatarPath)
	if err != nil {
		return
	}
	found = true
	return
}

func (db *dbManager) GetUserViaEmail(email string) (user models.User, found bool, err error) {

	row := db.dateBase.QueryRow(
		`SELECT * FROM public.users WHERE email = $1`, email)
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Salt, &user.Won, &user.Lost, &user.PlayTime, &user.Nickname,
		&user.AvatarPath)
	if err != nil {
		return
	}
	found = true
	return
}

func (db *dbManager) AddUser(user models.User) (err error) {

	_, err = db.dateBase.Exec(
		`INSERT INTO public.users (email, password, salt, won, lost, playtime, nickname, avatarpath)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.Email, user.Password, user.Salt, user.Won, user.Lost, user.PlayTime, user.Nickname, user.AvatarPath)
	return
}

func (db *dbManager) DeleteUser(email string) (err error) {

	_, err = db.dateBase.Exec(`DELETE FROM public.users WHERE email = $1`, email)
	return
}

func (db *dbManager) GetLenUsers() (len int, err error) {

	row := db.dateBase.QueryRow(`SELECT COUNT(*) FROM public.users`)
	err = row.Scan(&len)
	return
}

func (db *dbManager) GetUsers() (users []models.User, err error) {

	rows, err := db.dateBase.Query(
		`SELECT * FROM public.users ORDER BY won DESC`)
	if err != nil {
		return
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			return
		}
	}()

	var user models.User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Salt, &user.Won, &user.Lost, &user.PlayTime, &user.Nickname,
			&user.AvatarPath)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	return
}

func (db *dbManager) CleanerDBForTests() (err error) {
	_, err = db.dateBase.Exec(`TRUNCATE TABLE public.users`)
	return
}


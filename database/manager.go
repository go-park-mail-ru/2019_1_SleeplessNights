package database

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/lib/pq"
	"math/rand"
	"strconv"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1209qawsed"
	dbName   = "postgres"
)

//const (
//	host     = ""
//	port     = 0
//	user     = ""
//	password = ""
//	dbName   = ""
//)

const (
	CountOfPacks = 10
)

var db *dbManager

type dbManager struct {
	dataBase *sql.DB
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
		dataBase: dateBase,
	}

	return
}

func CloseConnection() (err error) {
	err = db.dataBase.Close()
	return
}

func GetInstance() *dbManager {
	return db
}

func (db *dbManager) GetUserViaID(userID uint) (user models.User, found bool, err error) {

	row := db.dataBase.QueryRow(
		`SELECT * FROM public.users WHERE id = $1`, userID)
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Salt, &user.Won, &user.Lost, &user.PlayTime, &user.Nickname,
		&user.AvatarPath)
	if err != nil {
		return
	}
	found = true
	return
}

func (db *dbManager) GetUserViaEmail(email string) (user models.User, err error) {

	row := db.dataBase.QueryRow(
		`SELECT * FROM public.users WHERE email = $1`, email)
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Salt, &user.Won, &user.Lost, &user.PlayTime, &user.Nickname,
		&user.AvatarPath)
	return
}

func (db *dbManager) AddUser(user models.User) (err error) {

	_, err = db.dataBase.Exec(
		`INSERT INTO public.users (email, password, salt, won, lost, playtime, nickname, avatarpath)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.Email, user.Password, user.Salt, user.Won, user.Lost, user.PlayTime, user.Nickname, user.AvatarPath)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}
	return
}

func (db *dbManager) UpdateUser(user models.User, email string) (err error) {

	_, err = db.dataBase.Exec(
		`UPDATE public.users 
			SET email = $1, password = $2, salt = $3, nickname = $4, avatarpath = $5
			WHERE email = $6`, user.Email, user.Password, user.Salt, user.Nickname, user.AvatarPath, email)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}
	return
}

func (db *dbManager) GetLenUsers() (len int, err error) {

	row := db.dataBase.QueryRow(`SELECT COUNT(*) FROM public.users`)
	err = row.Scan(&len)
	return
}

func (db *dbManager) GetUsers() (users []models.User, err error) {

	rows, err := db.dataBase.Query(
		`SELECT * FROM public.users ORDER BY won DESC`)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}

	var user models.User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Salt, &user.Won, &user.Lost, &user.PlayTime, &user.Nickname,
			&user.AvatarPath)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	err = rows.Close()
	if err != nil {
		return
	}
	return
}

func (db *dbManager) CleanerDBForTests() (err error) {
	_, err = db.dataBase.Exec(`TRUNCATE TABLE public.users RESTART IDENTITY`)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}
	return
}

func (db *dbManager) GetPackOfQuestions(theme string) (pack []models.Question, err error) {

	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(CountOfPacks)
	resultTheme := fmt.Sprintf(theme + strconv.Itoa(number))

	rows, err := db.dataBase.Query(
		`SELECT * FROM public.question WHERE pack_id = 
        (SELECT id FROM public.question_pack WHERE theme = $1)`,
		resultTheme)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}

	var question models.Question
	for rows.Next(){
		err = rows.Scan(&question.ID, &question.Answers, &question.Correct, &question.Text, &question.PackID)
		err = rows.Close()
		if err != nil {
			return
		}

		pack = append(pack, question)
	}

	err = rows.Close()
	if err != nil {
		return
	}
	return
}

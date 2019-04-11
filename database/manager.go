package database

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/lib/pq"
	"github.com/xlab/closer"
)

const (
	host     = ""
	port     = 0
	user     = ""
	password = ""
	dbName   = ""
)

var db *dbManager

type dbManager struct {
	dataBase *sql.DB
}

func init() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	dateBase, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal.Print(err.Error())
	}

	err = dateBase.Ping()
	if err != nil {
		logger.Fatal.Print(err.Error())
	}

	db = &dbManager{
		dataBase: dateBase,
	}

	closer.Bind(CloseConnection)
}

func CloseConnection() {
	err := db.dataBase.Close()
	if err != nil {
		logger.Fatal.Print(err.Error())
	}
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

func (db *dbManager) GetPacksOfQuestions(theme string) (packs map[string][]models.Question, err error) {

	rows, err := db.dataBase.Query(
		`SELECT * FROM public.question WHERE pack_id = 
        (SELECT DISTINCT ON (theme) id FROM public.question_pack ORDER BY theme, random() LIMIT 10 )`, theme)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}
	for rows.Next() {

		var question models.Question
		var theme string
		err = rows.Scan(&question.ID, &question.Answers, &question.Correct, &question.Text, &question.PackID, &theme)
		err = rows.Close()
		if err != nil {
			return
		}

		pack := packs[theme]
		pack = append(pack, question)
		packs[theme] = pack
	}

	err = rows.Close()
	if err != nil {
		return
	}
	return
}

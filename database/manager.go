package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/lib/pq"
	"github.com/xlab/closer"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1209qawsed"
	dbName   = "postgres"
)

const  (
	SQLNoRows   = "sql: no rows in result set"
	NoUserFound = "БД: Не был найден юзер"
)

var db *dbManager

type dbManager struct {
	dataBase *sql.DB
}

func init() {
	fmt.Println("Connection opened")
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

	closer.Bind(closeConnection)
}

func closeConnection() {
	err := db.dataBase.Close()
	if err != nil {
		logger.Fatal.Print(err.Error())
	}
	fmt.Println("Connection closed")
}

func GetInstance() *dbManager {
	return db
}

func (db *dbManager) GetUserViaID(userID uint) (user models.User, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	row := db.dataBase.QueryRow(
		`SELECT * FROM public.users WHERE id = $1`, userID)
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Salt, &user.Won, &user.Lost, &user.PlayTime, &user.Nickname,
		&user.AvatarPath)
	if err != nil && err.Error() == SQLNoRows {
		err = errors.New(NoUserFound)
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetUserViaEmail(email string) (user models.User, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	row := db.dataBase.QueryRow(
		`SELECT * FROM public.users WHERE email = $1`, email)
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Salt, &user.Won, &user.Lost, &user.PlayTime, &user.Nickname,
		&user.AvatarPath)
	if err != nil {
		err = errors.New(NoUserFound)
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) AddUser(user models.User) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(
		`INSERT INTO public.users (email, password, salt, won, lost, playtime, nickname, avatarpath)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.Email, user.Password, user.Salt, user.Won, user.Lost, user.PlayTime, user.Nickname, user.AvatarPath)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) UpdateUser(user models.User, userID uint) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(
		`UPDATE public.users 
			SET email = CASE
				WHEN $1 = '' THEN email ELSE $1 END,
			    nickname = CASE
				WHEN $2 = '' THEN nickname ELSE $2 END,
			    avatarpath = CASE
				WHEN $3 = '' THEN avatarpath ELSE $3 END
			WHERE id = $4 `, user.Email, user.Nickname, user.AvatarPath, userID)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetLenUsers() (len int, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	row := db.dataBase.QueryRow(`SELECT COUNT(*) FROM public.users`)
	err = row.Scan(&len)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetUsers() (users []models.User, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

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
	err = rows.Err()
	if err != nil {
		return
	}

	err = rows.Close()
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) CleanerDBForTests() (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(`TRUNCATE TABLE public.users, public.question, public.question_pack RESTART IDENTITY`)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetPacksOfQuestions(theme string) (packs map[string][]models.Question, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	rows, err := db.dataBase.Query(
		`SELECT * FROM public.question WHERE pack_id = 
        (SELECT DISTINCT ON (theme) id FROM public.question_pack ORDER BY theme, random() LIMIT 10 )`, theme)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}
	for rows.Next() {

		var question models.Question
		err = rows.Scan(&question.ID, &question.Answers, &question.Correct, &question.Text, &question.PackID, &question.Theme)

		pack := packs[question.Theme]
		pack = append(pack, question)
		packs[question.Theme] = pack
	}
	err = rows.Err()
	if err != nil {
		return
	}

	err = rows.Close()
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) AddQuestionPack(theme string) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(
		`INSERT INTO public.question_pack (theme) VALUES ($1)`, theme)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) AddQuestion(question models.Question) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(
		`INSERT INTO public.question (answers, correct, text, pack_id, pack_theme)
			  VALUES ($1, $2, $3, $4, $5)`, pq.Array(question.Answers), question.Correct, question.Text, question.PackID, question.Theme)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

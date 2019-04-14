package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/lib/pq"
	"github.com/xlab/closer"
	"os"
)

const (
	SQLNoRows   = "sql: no rows in result set"
	NoUserFound = "БД: Не был найден юзер"
)

var db *dbManager

var logger *log.Logger

func init () {
	logger = log.GetLogger("DB")
}

type dbManager struct {
	dataBase *sql.DB
}

type dbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func loadConfiguration(file string) (psqlInfo string) {
	configFile, err := os.Open(file)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	var config dbConfig
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = configFile.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	return
}

func init() {
	psqlInfo := loadConfiguration(os.Getenv("BASEPATH") + "/database/config.json")

	dateBase, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = dateBase.Ping()
	if err != nil {
		logger.Fatal(err.Error())
	}
	fmt.Println("DB connection opened")

	db = &dbManager{
		dataBase: dateBase,
	}

	closer.Bind(closeConnection)

}

func closeConnection() {
	err := db.dataBase.Close()
	if err != nil {
		logger.Fatal(err.Error())
	}
	fmt.Println("DB connection closed")
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
		logger.Error(_err.Error())
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
			SET nickname = CASE
				WHEN $1 = '' THEN nickname ELSE $1 END,
			    avatarpath = CASE
				WHEN $2 = '' THEN avatarpath ELSE $2 END
			WHERE id = $3`, user.Nickname, user.AvatarPath, userID)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error(_err.Error())
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
		logger.Error(_err.Error())
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
		logger.Error(_err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetPacksOfQuestions() (packs []models.Pack, err error) {

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
		`SELECT * FROM 
               (SELECT DISTINCT ON (theme) * FROM public.question_pack ORDER BY theme) AS qp
				ORDER BY random() LIMIT 10`)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error(_err.Error())
		return
	}

	var pack models.Pack
	for rows.Next() {

		err = rows.Scan(&pack.ID, &pack.Theme)
		if err != nil {
			return
		}

		packs = append(packs, pack)
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

func (db *dbManager) GetQuestions(ids []int) (questions []models.Question, err error) {

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

	//strOfIds := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")

	rows, err := db.dataBase.Query(
		`SELECT * FROM public.question WHERE pack_id = ANY ($1)`, pq.Array(ids))
	if _err, ok := err.(*pq.Error); ok {
		logger.Error(_err.Error())
		return
	}

	var question models.Question
	for rows.Next() {
		err = rows.Scan(&question.ID, pq.Array(&question.Answers), &question.Correct, &question.Text, &question.PackID)
		if err != nil {
			return
		}

		questions = append(questions, question)
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
		logger.Error(_err.Error())
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
		`INSERT INTO public.question (answers, correct, text, pack_id)
			  VALUES ($1, $2, $3, $4)`, pq.Array(question.Answers), question.Correct, question.Text, question.PackID)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error(_err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	txOK = true
	return
}

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/models"
	"github.com/lib/pq"
	"github.com/xlab/closer"
	"os"
)

const (
	SQLNoRows       = "sql: no rows in result set"
	NoUserFound     = "БД: Не был найден юзер"
	UniqueViolation = "pq: duplicate key value violates unique constraint \"users_email_ui\""
)

var db *dbManager

var logger *log.Logger

func init() {
	logger = log.GetLogger("DB")
}

type dbManager struct {
	dataBase *sql.DB
}

type dbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user_manager"`
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
	psqlInfo = fmt.Sprintf("host=%s port=%d user_manager=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	return
}

func init() {
	psqlInfo := loadConfiguration(os.Getenv("BASEPATH") + "/main_microservice/database/microservices.json")

	dataBase, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = dataBase.Ping()
	if err != nil {
		logger.Fatal(err.Error())
	}
	fmt.Println("DB connection opened")

	db = &dbManager{
		dataBase: dataBase,
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

func (db *dbManager) CleanerDBForTests() (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(`TRUNCATE TABLE public.question, public.question_pack RESTART IDENTITY`)
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

func (db *dbManager) GetPacksOfQuestions(n int) (packs []models.Pack, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	rows, err := db.dataBase.Query(
		`SELECT * FROM 
               (SELECT DISTINCT ON (theme) * FROM public.question_pack ORDER BY theme) AS qp
				ORDER BY random() LIMIT $1`, n)
	if err != nil {
		logger.Error(err.Error())
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
			_ = tx.Rollback()
		}
	}()

	rows, err := db.dataBase.Query(
		`SELECT * FROM public.question WHERE pack_id = ANY ($1)`, pq.Array(ids))
	if err != nil {
		logger.Error(err.Error())
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
			_ = tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(
		`INSERT INTO public.question_pack (theme) VALUES ($1)`, theme)
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

func (db *dbManager) AddQuestion(question models.Question) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	_, err = db.dataBase.Exec(
		`INSERT INTO public.question (answers, correct, text, pack_id)
			  VALUES ($1, $2, $3, $4)`, pq.Array(question.Answers), question.Correct, question.Text, question.PackID)
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

package database

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database/models"
	"github.com/jackc/pgx"
	"github.com/lib/pq"
	"github.com/xlab/closer"
	"os"
	"time"

	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
)

const (
	maxConnections = 3
	acquireTimeout = 3 * time.Second
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("DB")
}

type dbConfig struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func loadConfiguration(file string) (pgxConfig pgx.ConnConfig) {
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

	pgxConfig.Host = config.Host
	pgxConfig.User = config.User
	pgxConfig.Password = config.Password
	pgxConfig.Database = config.DBName
	pgxConfig.Port = config.Port

	return
}

var db *dbManager

type dbManager struct {
	dataBase *pgx.ConnPool
}

func init() {
	//TODO check config loading
	pgxConfig := loadConfiguration(os.Getenv("BASEPATH") + "/game_microservice/database/config.json")
	pgxConnPoolConfig := pgx.ConnPoolConfig{ConnConfig: pgxConfig, MaxConnections: maxConnections, AcquireTimeout: acquireTimeout}

	dataBase, err := pgx.NewConnPool(pgxConnPoolConfig)
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
	db.dataBase.Close()
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

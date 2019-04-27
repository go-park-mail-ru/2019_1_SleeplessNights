package database

import (
	"encoding/json"
	"fmt"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/jackc/pgx"
	"github.com/xlab/closer"
	"os"
	"time"
)

const (
	maxConnections = 3
	acquireTimeout = 3 * time.Second
)

var db *dbManager

var logger *log.Logger

func init() {
	logger = log.GetLogger("DB")
}

type dbManager struct {
	dataBase *pgx.ConnPool
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

func init() {

	pgxConfig := loadConfiguration(os.Getenv("BASEPATH") + "/chat_microservice/database/config.json")
	pgxConnPoolConfig := pgx.ConnPoolConfig{ConnConfig: pgxConfig, MaxConnections: maxConnections, AcquireTimeout: acquireTimeout}

	dateBase, err := pgx.NewConnPool(pgxConnPoolConfig)
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
	db.dataBase.Close()
	fmt.Println("DB connection closed")
}

func GetInstance() *dbManager {
	return db
}

func (db *dbManager) CreateMessage() (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetMessages() (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	txOK := false
	defer func() {
		if !txOK {
			_ = tx.Rollback()
		}
	}()

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	txOK = true
	return
}

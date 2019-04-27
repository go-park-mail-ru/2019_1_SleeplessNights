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

func (db *dbManager) PostMessage(userId uint64, roomId uint64, payload []byte) (err error) {

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

	_, err = db.dataBase.Exec(`SELECT * FROM func_post_message ($1, $2, $3)`,
		userId, payload, roomId)
	if err != nil{
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	txOK = true
	return
}

func (db *dbManager) GetMessages(roomId uint64, since uint64, limit uint64) (payload []byte, err error) {

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

	row := db.dataBase.QueryRow(`SELECT * FROM func_get_messages ($1, $2, $3)`,
		roomId, since, limit)
	err = row.Scan(&payload)
	if err != nil{
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	txOK = true
	return
}

func (db *dbManager) UpdateUser(uid uint64, nickname string, avatarPath string) (id uint64, err error) {

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

	row := db.dataBase.QueryRow(`SELECT * FROM func_update_user ($1, $2, $3)`,
		uid, nickname, avatarPath)
	err = row.Scan(&id)
	if err != nil{
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	txOK = true
	return
}

func (db *dbManager) CreateRoom(users []uint64) (id uint64, err error) {

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

	row := db.dataBase.QueryRow(`SELECT * FROM func_create_room ($1)`, users)
	err = row.Scan(&id)
	if err != nil{
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	txOK = true
	return
}
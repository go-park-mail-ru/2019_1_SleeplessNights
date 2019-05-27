package database

import (
	"encoding/json"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx"
	"github.com/xlab/closer"
	"os"
	"time"
)

const (
	maxConnections = 3
	acquireTimeout = 3 * time.Second
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("DataBase")
	logger.SetLogLevel(logrus.TraceLevel)
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
		logger.Errorf("Failed to open file: %v", err.Error())
		return
	}
	var config dbConfig
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		logger.Errorf("Failed to decode config: %v", err.Error())
		return
	}
	err = configFile.Close()
	if err != nil {
		logger.Errorf("Failed to close file: %v", err.Error())
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

	pgxConfig := loadConfiguration(os.Getenv("BASEPATH") + "/user_microservice/database/config.json")
	pgxConnPoolConfig := pgx.ConnPoolConfig{ConnConfig: pgxConfig, MaxConnections: maxConnections, AcquireTimeout: acquireTimeout}

	dataBase, err := pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		logger.Fatalf("Failed to get conn pool: %v", err.Error())
	}

	logger.Info("DB connection opened")

	db = &dbManager{
		dataBase: dataBase,
	}

	closer.Bind(closeConnection)
}
func closeConnection() {
	db.dataBase.Close()
	logger.Info("DB connection closed")
}

func GetInstance() *dbManager {
	return db
}

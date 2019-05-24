package database

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"time"

	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
)

var (
	maxConnections = config.GetInt("game_ms.pkg.database.max_connections")
	acquireTimeout = config.GetDuration("game_ms.pkg.database.acquire_timeout", 3*time.Second)
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("DB")
	logger.SetLogLevel(logrus.Level(config.GetInt("game_ms.log_level")))
}

func loadConfiguration() (pgxConfig pgx.ConnConfig) {
	pgxConfig.Port = uint16(config.GetInt("postgres.port"))
	pgxConfig.Host = config.GetString("postgres.host")
	pgxConfig.Database = config.GetString("postgres.db_name")
	pgxConfig.User = config.GetString("postgres.user")
	pgxConfig.Password = config.GetString("postgres.password")
	return
}

var db *dbManager

type dbManager struct {
	dataBase *pgx.ConnPool
}

func init() {
	//TODO check config loading
	pgxConfig := loadConfiguration()
	pgxConnPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     pgxConfig,
		MaxConnections: maxConnections,
		AcquireTimeout: acquireTimeout,
	}

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

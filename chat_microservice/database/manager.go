package database

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"time"
)

var (
	maxConnections = config.GetInt("user_ms.pkg.database.max_connections")
	acquireTimeout time.Duration
)

func init() {
	var err error
	acquireTimeout, err = time.ParseDuration(config.GetString("user_ms.pkg.database.acquire_timeout"))
	if err != nil {
		acquireTimeout = 3 * time.Second
	}
}

var db *dbManager

var logger *log.Logger

func init() {
	logger = log.GetLogger("DataBase")
	logger.SetLogLevel(logrus.Level(config.GetInt("chat_ms.log_level")))
}

type dbManager struct {
	dataBase *pgx.ConnPool
}

func loadConfiguration() (pgxConfig pgx.ConnConfig) {
	pgxConfig.Port = uint16(config.GetInt("postgres.port"))
	pgxConfig.Host = config.GetString("postgres.host")
	pgxConfig.Database = config.GetString("postgres.db_name")
	pgxConfig.User = config.GetString("postgres.user")
	pgxConfig.Password = config.GetString("postgres.password")
	return
}

func init() {

	pgxConfig := loadConfiguration()
	pgxConnPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgxConfig,
		MaxConnections: maxConnections,
		AcquireTimeout: acquireTimeout,
	}

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

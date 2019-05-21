package database

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/jackc/pgx"
	"github.com/xlab/closer"
)

var (
	maxConnections = config.GetInt("user_ms.pkg.database.max_connections")
	acquireTimeout time.Duration
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("DataBase")
	logger.SetLogLevel(logrus.Level(config.GetInt("user_ms.log_level")))
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

	pgxConfig := loadConfiguration()
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

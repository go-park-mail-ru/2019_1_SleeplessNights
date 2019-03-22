package database

import (
	"github.com/jackc/pgx"
	"log"
)

type dbManager struct {
	connections *pgx.ConnPool
}

var manager *dbManager

func init() {
	pool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		log.Fatal(err)
	}

	manager = &dbManager {
		connections: pool,
	}
}

func GetInstance() *dbManager {
	return manager
}

func (manager *dbManager)Exec(sql SQL)(err error) {
	return nil
}

func (manager *dbManager)Query(sql SQL, map[string]re)(tuples []interface{}, err error) {
	return nil, nil
}

func (manager *dbManager)QuerySingle(sql SQL)(tuple interface{}, err error) {
	return nil, nil
}
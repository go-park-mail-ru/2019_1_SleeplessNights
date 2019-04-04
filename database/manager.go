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

func (manager *dbManager)Exec(sql string, args ...interface{})(err error) {
	tx, err := manager.connections.Begin()
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-opx
	defer tx.Rollback()

	_, err = tx.Exec(sql, args)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return
}

func (manager *dbManager)Query(sql string, args ...interface{})(rows *pgx.Rows, err error) {
	tx, err := manager.connections.Begin()
	if err != nil {
		return nil, err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()
	rows, err = tx.Query(sql, args)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	return
}

func (manager *dbManager)QueryRow(sql string, args ...interface{})(row *pgx.Row, err error) {
	tx, err := manager.connections.Begin()
	if err != nil {
		return nil, err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()

	row = tx.QueryRow(sql, args[0], args[1], args[2])//TODO FIX
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	return
}

package service

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/jackc/pgx"
)

func clear()(err error) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM func_service_clear()`)
	if err == nil {
		err = tx.Commit()
	}
	return
}

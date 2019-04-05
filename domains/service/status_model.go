package service

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
)

func status()(code int, response interface{}) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	var status = struct {
		User   int `json:"user"`
		Forum  int `json:"forum"`
		Thread int `json:"thread"`
		Post   int `json:"post"`
	}{}

	row := tx.QueryRow(`SELECT * FROM func_service_status()`)
	err = row.Scan(&status.User, &status.Forum, &status.Thread, &status.Post)
	if err == nil {
		code = 200
		response = &status
		err = tx.Commit()
		if err != nil {
			return responses.InternalError("Error while committing transaction: " + err.Error())
		}
	} else {
		return responses.InternalError("Database returned unexpected error: " + err.Error())
	}
	return
}

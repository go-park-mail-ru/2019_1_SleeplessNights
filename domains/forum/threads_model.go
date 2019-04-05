package forum

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
	"time"
)

func threads(slug string, limit int32, since time.Time, desc bool)(code int, response interface{}) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	var sincePtr *time.Time
	if since.IsZero() {
		sincePtr = nil
	} else {
		sincePtr = &since
	}

	var rows *pgx.Rows
	rows, err = tx.Query(`SELECT * FROM func_forum_threads($1, $2, $3, $4)`, slug, sincePtr, desc, limit)
	defer rows.Close()

	if err == nil {
		var threads []responses.Thread
		for rows.Next() {
			var thread responses.Thread
			err = rows.Scan(&thread.IsNew, &thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
			if err != nil {
				return responses.InternalError("Error while scanning row: " + err.Error())
			}
			threads = append(threads, thread)
		}
		if rows.Err() != nil {
			return responses.InternalError("Error returned by rows: " + err.Error())
		}

		code = 200
		response = &threads
		err = tx.Commit()
		if err != nil {
			return responses.InternalError("Error while committing transaction: " + err.Error())
		}
	} else {
		switch err.Error() {
		case "ERROR: no_data_found (SQLSTATE P0002)":
			code = 404
			var msg responses.Error
			msg.Message = "Can't find forum"
			response = &msg
		default:
			return responses.InternalError("Database returned unexpected error: " + err.Error())
		}
	}
	return
}
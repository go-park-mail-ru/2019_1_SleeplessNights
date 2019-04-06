package forum

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
	"time"
)

func createThread(slug, title, author, message, threadSlug string, created time.Time)(code int, response interface{}) {
	/*removeSpaces := func(r rune) rune{
		if r == ' ' {
			return '-'
		} else {
			return r
		}
	}
	threadSlug := strings.Map(removeSpaces, strings.ToLower(title))*/

	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()
	if created.IsZero() {
		created = time.Now()
	}
	row := tx.QueryRow(`SELECT * FROM func_forum_create_thread($1, $2, $3, $4, $5, $6)`,
		slug, threadSlug, title, author, message, created)
	var thread responses.Thread
	err = row.Scan(&thread.IsNew, &thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err == nil {
		if thread.IsNew {
			code = 201
		} else {
			code = 409
		}
		response = &thread
		err = tx.Commit()
		if err != nil {
			return responses.InternalError("Error while committing transaction: " + err.Error())
		}
	} else {
		switch err.Error() {
		case "ERROR: no_data_found (SQLSTATE P0002)":
			code = 404
			var msg responses.Error
			msg.Message = "Can't find forum or user"
			response = &msg
		default:
			return responses.InternalError("Database returned unexpected error: " + err.Error())
		}
	}
	return
}

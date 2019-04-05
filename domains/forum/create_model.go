package forum

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
)

func create(title, user, slug string)(code int, response interface{}) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_forum_create($1, $2, $3)`, title, user, slug)
	var forum responses.Forum
	err = row.Scan(&forum.IsNew, &forum.ForumTitle, &forum.UserNickname, &forum.ForumSlug, &forum.PostsCount, &forum.ThreadsCount)
	if err == nil {
		if forum.IsNew {
			code = 201
		} else {
			code = 409
		}
		response = &forum
		err = tx.Commit()
		if err != nil {
			return responses.InternalError("Error while committing transaction: " + err.Error())
		}
	} else {
		switch err.Error() {
		case "ERROR: no_data_found (SQLSTATE P0002)":
			code = 404
			var msg responses.Error
			msg.Message = "Can't find user"
			response = &msg
		default:
			return responses.InternalError("Database returned unexpected error: " + err.Error())
		}
	}
	return
}

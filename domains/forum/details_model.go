package forum

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
)

func details(slug string)(code int, response interface{}) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_forum_details($1)`, slug)
	var forum responses.Forum
	err = row.Scan(&forum.ForumTitle, &forum.UserNickname, &forum.ForumSlug, &forum.PostsCount, &forum.ThreadsCount, &forum.IsNew)
	if err == nil {
		code = 200
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
			msg.Message = "Can't find forum"
			response = &msg
		default:
			return responses.InternalError("Database returned unexpected error: " + err.Error())
		}
	}
	return
}
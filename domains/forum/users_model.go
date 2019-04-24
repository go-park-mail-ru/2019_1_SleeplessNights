package forum

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
)

func users(slug string, limit int32, since string, desc bool)(code int, response interface{}) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	var sincePtr *string
	if since == "" {
		sincePtr = nil
	} else {
		sincePtr = &since
	}

	var rows *pgx.Rows
	rows, err = tx.Query(`SELECT * FROM func_forum_users($1, $2, $3, $4)`, slug, sincePtr, desc, limit)
	defer rows.Close()

	if err == nil {
		var users []responses.User
		for rows.Next() {
			var user responses.User
			err = rows.Scan(&user.IsNew, &user.Nickname, &user.Fullname, &user.About, &user.Email)
			if err != nil {
				return responses.InternalError("Error while scanning row: " + err.Error())
			}
			users = append(users, user)
		}
		err = rows.Err()
		if err != nil {
			switch err.Error() {
			case "ERROR: no_data_found (SQLSTATE P0002)":
				code = 404
				var msg responses.Error
				msg.Message = "Can't find forum"
				response = &msg
			default:
				return responses.InternalError("Database returned unexpected error: " + err.Error())
			}
			return
		}

		code = 200
		response = &users
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

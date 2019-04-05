package user

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
)

func editProfile(nickname, fullname, about, email string)(code int, response interface{}) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_user_change_profile($1, $2, $3, $4)`, nickname, fullname, about, email)
	var user responses.User
	err = row.Scan(&user.IsNew, &user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err == nil {
		code = 200
		response = &user
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
		case "ERROR: unique_violation (SQLSTATE 23505)":
			code = 409
			var msg responses.Error
			msg.Message = "User with this email already exist"
			response = &msg
		default:
			return responses.InternalError("Database returned unexpected error: " + err.Error())
		}
	}
	return
}

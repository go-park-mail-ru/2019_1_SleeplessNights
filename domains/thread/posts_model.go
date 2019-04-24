package thread

import (
	"fmt"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
	"strconv"
)

const(
	flatSort = "flat"
	treeSort = "tree"
	parentTreeSort = "parent_tree"
)

func posts(slugOrID string, limit int32, since int64, sort string, desc bool)(code int, response interface{}) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	id, err := strconv.ParseInt(slugOrID, 10, 64)
	if err != nil {
		id, err = getIdBySlug(slugOrID, tx)
		if err != nil {
			switch err.Error() {
			case "ERROR: no_data_found (SQLSTATE P0002)":
				code = 404
				var msg responses.Error
				msg.Message = "Can't find thread"
				response = &msg
				return
			default:
				return responses.InternalError("Database returned unexpected error: " + err.Error())
			}
		}
	}

	var rows *pgx.Rows
	switch sort {
	case flatSort:
		rows, err = tx.Query(`SELECT * FROM func_thread_posts_flat($1, $2, $3, $4)`, id, since, desc, limit)
	case treeSort:
		rows, err = tx.Query(`SELECT * FROM func_thread_posts_tree($1, $2, $3, $4)`, id, since, desc, limit)
	case parentTreeSort:
		fmt.Println("PARENT TREE:",id, since, desc, limit)
		rows, err = tx.Query(`SELECT * FROM func_thread_posts_parent_tree($1, $2, $3, $4)`, id, since, desc, limit)
	}
	defer rows.Close()

	if err == nil {
		var posts []responses.Post
		var post responses.Post
		for rows.Next() {
			err = rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created, &post.IsNew)
			if err != nil {
				return responses.InternalError("Error while scanning row: " + err.Error())
			}
			posts = append(posts, post)
		}
		err = rows.Err()
		if err != nil {
			switch err.Error() {
			case "ERROR: no_data_found (SQLSTATE P0002)":
				code = 404
				var msg responses.Error
				msg.Message = "Can't find thread whits this id"
				response = &msg
				return
			case "ERROR: integrity_constraint_violation (SQLSTATE 23000)":
				code = 409
				var msg responses.Error
				msg.Message = "Can't find post in this thread"
				response = &msg
				return
			default:
				return responses.InternalError("Database returned unexpected error: " + err.Error())
			}
		}

		code = 200
		response = &posts
		err = tx.Commit()
		if err != nil {
			return responses.InternalError("Error while committing transaction: " + err.Error())
		}
	} else {
		return responses.InternalError("Error returned with rows: " + err.Error())
	}
	return
}

package post

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
)

const (
	keyUser =   "user"
	keyForum =  "forum"
	keyThread = "thread"
)

func details(id int64, related []string)(code int, response interface{}) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	tx, err := conn.Begin()
	if err != nil {
		return responses.InternalError("Error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	postInfo := struct {
		Post   responses.Post   `json:"post"`
		Author *responses.User   `json:"author,omitempty"`
		Thread *responses.Thread `json:"thread,omitempty"`
		Forum  *responses.Forum  `json:"forum,omitempty"`
	}{Author: nil, Thread: nil, Forum: nil}

	row := tx.QueryRow(`SELECT * FROM func_post_details($1)`, id)
	err = row.Scan(&postInfo.Post.ID, &postInfo.Post.Parent, &postInfo.Post.Author, &postInfo.Post.Message,
		&postInfo.Post.IsEdited, &postInfo.Post.Forum, &postInfo.Post.Thread, &postInfo.Post.Created, &postInfo.Post.IsNew)
	if err != nil{
		switch err.Error() {
		case "ERROR: no_data_found (SQLSTATE P0002)":
			code = 404
			var msg responses.Error
			msg.Message = "Can't find post"
			response = &msg
			return
		default:
			return responses.InternalError("Database returned unexpected error: " + err.Error())
		}
	}
	for _, key := range related {
		switch key {
		case "":
		case keyUser:
			postInfo.Author = &responses.User{}
			row := tx.QueryRow(`SELECT * FROM func_user_details($1)`, postInfo.Post.Author)
			err = row.Scan(&postInfo.Author.IsNew, &postInfo.Author.Nickname, &postInfo.Author.Fullname, &postInfo.Author.About, &postInfo.Author.Email)
			if err != nil{
				switch err.Error() {
				case "ERROR: no_data_found (SQLSTATE P0002)":
					code = 404
					var msg responses.Error
					msg.Message = "Can't find user"
					response = &msg
					return
				default:
					return responses.InternalError("Database returned unexpected error: " + err.Error())
				}
			}
		case keyForum:
			postInfo.Forum = &responses.Forum{}
			row := tx.QueryRow(`SELECT * FROM func_forum_details($1)`, postInfo.Post.Forum)
			err = row.Scan(&postInfo.Forum.ForumTitle, &postInfo.Forum.UserNickname, &postInfo.Forum.ForumSlug, &postInfo.Forum.PostsCount, &postInfo.Forum.ThreadsCount, &postInfo.Forum.IsNew)
			if err != nil{
				switch err.Error() {
				case "ERROR: no_data_found (SQLSTATE P0002)":
					code = 404
					var msg responses.Error
					msg.Message = "Can't find forum"
					response = &msg
					return
				default:
					return responses.InternalError("Database returned unexpected error: " + err.Error())
				}
			}
		case keyThread:
			postInfo.Thread = &responses.Thread{}
			row := tx.QueryRow(`SELECT * FROM func_thread_details($1)`, postInfo.Post.Thread)
			err = row.Scan(&postInfo.Thread.IsNew, &postInfo.Thread.ID, &postInfo.Thread.Title, &postInfo.Thread.Author,
				&postInfo.Thread.Forum, &postInfo.Thread.Message, &postInfo.Thread.Votes, &postInfo.Thread.Slug, &postInfo.Thread.Created)
			if err != nil{
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
		default:
			return 400, nil
		}
	}

	code = 200
	response = &postInfo
	err = tx.Commit()
	if err != nil {
		return responses.InternalError("Error while committing transaction: " + err.Error())
	}
	return
}

package forum

import (
	"fmt"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database/query"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
	"github.com/jackc/pgx"
)

func Create(slug, title, userNick string) (int, responses.ResponseModel) {
	conn, err := pgx.Connect(database.ConnConfig)
	defer conn.Close()
	if err != nil {
		fmt.Println("Error while connecting to DB:", err)
		data := responses.Error{Message:"Can't connect to the DB: " + err.Error()}
		return 500, &data
	}

	userRow := conn.QueryRow(query.FindUserIdByNicknameQuery, userNick)
	var userID uint64
	err = userRow.Scan(&userID)
	switch err {
	case nil: //Do nothing
	case pgx.ErrNoRows:
		fmt.Println("User not found")
		data := responses.Error{Message:"Can't find user with nickname " + userNick}
		return 404, &data
	default:
		fmt.Println("Error while searching for user:", err)
		data := responses.Error{Message:"Error while searching for user: " + err.Error()}
		return 500, &data
	}

	forumRow := conn.QueryRow(query.FindForumBySlugQuery, slug)
	var forum responses.Forum
	err = forumRow.Scan(&forum.PostsCount,
						&forum.ForumSlug,
						&forum.ThreadsCount,
						&forum.ForumTitle,
						&forum.UserNickname)
	switch err {
	case nil:
		fmt.Println("Forum found")
		return 409, &forum
	case pgx.ErrNoRows: //Do nothing
	default:
		fmt.Println("Error while searching for forum:", err)
		data := responses.Error{Message:"Error while searching for forum: " + err.Error()}
		return 500, &data
	}

	_, err = conn.Exec(query.CreateNewForum, slug, title, userID)
	if err != nil {
		fmt.Println("Error while creating new forum:", err)
		data := responses.Error{Message:"Error while creating new forum: " + err.Error()}
		return 500, &data
	}

	forumRow = conn.QueryRow(query.FindForumBySlugQuery, slug)
	err = forumRow.Scan(&forum.PostsCount,
		&forum.ForumSlug,
		&forum.ThreadsCount,
		&forum.ForumTitle,
		&forum.UserNickname)
	if err != nil {
		fmt.Println("Error while searching for created forum:", err)
		data := responses.Error{Message:"Error while searching for created forum: " + err.Error()}
		return 500, &data
	}

	fmt.Println("Forum created and found!")
	return 201, &forum
}
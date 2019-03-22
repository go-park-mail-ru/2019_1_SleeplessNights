package forum

import (
	"fmt"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/database/query"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/responses"
)

func Create(slug, title, userSlug string) (int, responses.ResponseModel,error) {
	dbManager := database.GetInstance()

	findUserQuery := query.FindUserIdBySlug{Slug: userSlug}
	userID, err := dbManager.QuerySingle(&findUserQuery)
	if err != nil {
		fmt.Println("Error while searching for user:", err)
		data := responses.Error{Message:"Error while searching for user: " + err.Error()}
		return 500, &data, err
	}
	if userID == nil {
		fmt.Println("User not found")
		data := responses.Error{Message:"Can't find user with slug " + userSlug}
		return 404, &data, nil
	}

	findForumQuery := query.FindForumBySlug{Slug: slug}
	forumID, err := dbManager.QuerySingle(&findForumQuery)
	if err != nil {
		fmt.Println("Error while searching for forum:", err)
		data := responses.Error{Message:"Error while searching for user: " + err.Error()}
		return 500, &data, err
	}
	if forumID != nil {

	}
}
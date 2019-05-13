package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"testing"
)

func TestGetUsersSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	firstUser := services.User{
		Id:         1,
		Email:      "test1@test.com",
		Nickname:   "test1",
		AvatarPath: "default_avatar1.jpg",
	}

	secondUser := services.User{
		Id:         2,
		Email:      "test2@test.com",
		Nickname:   "test2",
		AvatarPath: "default_avatar2.jpg",
	}

	_, err = database.GetInstance().AddUser(firstUser.Email, firstUser.Nickname, firstUser.AvatarPath, []byte{}, []byte{})
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = database.GetInstance().AddUser(secondUser.Email, secondUser.Nickname, secondUser.AvatarPath, []byte{}, []byte{})
	if err != nil {
		t.Error(err.Error())
		return
	}

	page := services.PageData{
		Since: 0,
		Limit: 100,
	}

	users, err := database.GetInstance().GetUsers(&page)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}

	defer func (){
		if recovered := recover(); recovered != nil {
			t.Errorf("DB returned wrong users or something alse")
		}
	}()

	if users.Leaders[0].User.Id != firstUser.Id {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			users.Leaders[0].User.Id, firstUser.Id)
	}
	if users.Leaders[0].User.Email != firstUser.Email {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			users.Leaders[0].User.Email, firstUser.Email)
	}
	if users.Leaders[0].User.Nickname != firstUser.Nickname {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			users.Leaders[0].User.Nickname, firstUser.Nickname)
	}
	if users.Leaders[0].User.AvatarPath != firstUser.AvatarPath {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			users.Leaders[0].User.AvatarPath, firstUser.AvatarPath)
	}
	if users.Leaders[1].User.Id != secondUser.Id {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			users.Leaders[1].User.Id, secondUser.Id)
	}
	if users.Leaders[1].User.Email != secondUser.Email {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			users.Leaders[1].User.Email, secondUser.Email)
	}
	if users.Leaders[1].User.Nickname != secondUser.Nickname {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			users.Leaders[1].User.Nickname, secondUser.Nickname)
	}
	if users.Leaders[1].User.AvatarPath != secondUser.AvatarPath {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			users.Leaders[1].User.AvatarPath, secondUser.AvatarPath)
	}
}
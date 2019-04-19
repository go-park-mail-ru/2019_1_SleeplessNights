package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"testing"
)

func TestGetUserViaIDSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUser := models.User{
		ID:         1,
		Email:      "first@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	err = database.GetInstance().AddUser(oldUser)
	if err != nil {
		t.Error(err.Error())
		return
	}

	newUser, err := database.GetInstance().GetUserViaID(oldUser.ID)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if newUser.ID != oldUser.ID || newUser.Email != oldUser.Email {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			newUser.ID, newUser.Email, oldUser.ID, oldUser.Email)
	}
}

func TestGetUserViaIDUnsuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	var userID uint = 1
	expected := database.NoUserFound

	_, err = database.GetInstance().GetUserViaID(userID)
	if err == nil {
		t.Errorf("DB didn't return error")
	} else if err.Error() != expected {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v\n",
			err.Error(), expected)
	}
}

func TestGetUserViaEmailSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUser := models.User{
		ID:         1,
		Email:      "first@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	err = database.GetInstance().AddUser(oldUser)
	if err != nil {
		t.Error(err.Error())
		return
	}

	newUser, err := database.GetInstance().GetUserViaEmail(oldUser.Email)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if newUser.Email != oldUser.Email || newUser.ID != oldUser.ID {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			newUser.Email, newUser.ID, oldUser.Email, oldUser.ID)
	}
}

func TestGetUserViaEmailUnsuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	var userEmail = "test@test.com"
	expected := database.NoUserFound

	_, err = database.GetInstance().GetUserViaEmail(userEmail)
	if err == nil {
		t.Errorf("DB didn't return error")
	} else if err.Error() != expected {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v\n",
			err.Error(), expected)
	}
}

func TestAddUserSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUser := models.User{
		ID:         1,
		Email:      "test@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	err = database.GetInstance().AddUser(oldUser)
	if err != nil {
		t.Error(err.Error())
		return
	}

	newUser, err := database.GetInstance().GetUserViaEmail(oldUser.Email)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if newUser.Email != oldUser.Email || newUser.ID != oldUser.ID {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			newUser.Email, newUser.ID, oldUser.Email, oldUser.ID)
	}
}

func TestAddUserUnsuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	firstUser := models.User{
		Email:      "first@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	secondUser := models.User{
		Email:      "first@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	err = database.GetInstance().AddUser(firstUser)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	expected := database.UniqueViolation

	err = database.GetInstance().AddUser(secondUser)
	if err == nil {
		t.Errorf("DB didn't return error")
	} else if err.Error() != expected {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v\n",
			err.Error(), expected)
	}
}

func TestUpdateUserSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUser := models.User{
		Email:      "test@test.com",
		Nickname:   "old",
		AvatarPath: "old_default_avatar.jpg",
	}

	err = database.GetInstance().AddUser(oldUser)
	if err != nil {
		t.Error(err.Error())
		return
	}

	oldUser = models.User{
		ID:         1,
		Nickname:   "new",
		AvatarPath: "new_default_avatar.jpg",
	}

	err = database.GetInstance().UpdateUser(oldUser, oldUser.ID)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}

	newUser, err := database.GetInstance().GetUserViaID(oldUser.ID)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if newUser.Nickname != oldUser.Nickname || newUser.AvatarPath != oldUser.AvatarPath {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			newUser.Nickname, oldUser.Nickname, newUser.AvatarPath, oldUser.AvatarPath)
	}
}

//func TestUpdateUserUnsuccessful(t *testing.T) {
//
//	err := database.GetInstance().CleanerDBForTests()
//	if err != nil {
//	t.Errorf(err.Error())
//	}
//
//	oldUser := models.User{
//		ID:         2000,
//		Nickname:   "temp",
//		AvatarPath: "new_default_avatar.jpg",
//	}
//
//	err = database.GetInstance().UpdateUser(oldUser, oldUser.ID)
//	if err == nil {
//		t.Errorf("DB didn't return error")
//	}
//
//	err = database.GetInstance().CleanerDBForTests()
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//}

func TestGetLenUsers(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	user := models.User{
		Email:      "test@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	err = database.GetInstance().AddUser(user)
	if err != nil {
		t.Error(err.Error())
		return
	}

	newLen, err := database.GetInstance().GetLenUsers()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if newLen != 1 {
		t.Errorf("DB returned wrong length:\ngot %v\nwant %v\n",
			newLen, 1)
	}
}

func TestGetUsers(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUsers := []models.User{
		{
			Email:      "first@test.com",
			Nickname:   "first",
			Won:        10,
			AvatarPath: "first_default_avatar.jpg",
		},
		{
			Email:      "second@test.com",
			Won:        1,
			Nickname:   "second",
			AvatarPath: "second_default_avatar.jpg",
		},
	}

	for _, user := range oldUsers {
		err := database.GetInstance().AddUser(user)
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	newUsers, err := database.GetInstance().GetUsers()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	for i, _ := range newUsers {
		if newUsers[i].Email != oldUsers[i].Email || newUsers[i].Nickname != oldUsers[i].Nickname || newUsers[i].AvatarPath != oldUsers[i].AvatarPath {
			t.Errorf("DB returned wrong user:\ngot %v, %v, %v\nwant %v, %v, %v\n",
				newUsers[i].Email, newUsers[i].Nickname, newUsers[i].AvatarPath,
				oldUsers[i].Email, oldUsers[i].Nickname, oldUsers[i].AvatarPath)
		}
	}
}

func TestCleanerDBForTests(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	user := models.User{
		Email:      "test@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	err = database.GetInstance().AddUser(user)
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	length, err := database.GetInstance().GetLenUsers()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if length != 0 {
		t.Errorf("DB didn't cleaned up")
	}
}

func TestPacksOfQuestions(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	//themes := []string{
	//	"математика",
	//	"информатика",
	//	"химия",
	//	"биология",
	//	"физика",
	//	"культура",
	//	"история",
	//	"языки",
	//	"философия",
	//	"мемология",
	//}
	//
	//for _, theme := range themes {
	//	err = database.GetInstance().AddQuestionPack(theme)
	//	if err != nil {
	//		t.Errorf("DB returned error: %v", err.Error())
	//	}
	//}

	faker.CreateFakePacks()

	packs, err := database.GetInstance().GetPacksOfQuestions(10)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	count := len(packs)
	if count != 10 {
		t.Errorf("DB return wrong count of packs: %v", count)
	}
}

func TestGetQuestions(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	//for i := uint(1); i <= uint(4); i++  {
	//	question := models.Question{
	//		Answers: []string{},
	//		Correct: 2,
	//		Text:    "",
	//		PackID:  i,
	//	}
	//
	//	err = database.GetInstance().AddQuestion(question)
	//	if err != nil {
	//		t.Errorf("DB returned error: %v", err.Error())
	//	}
	//}

	faker.CreateFakePacks()

	questions, err := database.GetInstance().GetQuestions([]int{1, 2, 3, 4})
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	cont := len(questions)
	if cont != 40 {
		t.Errorf("DB return wrong count of questions")
	}
}

func TestAddQuestionPack(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	err = database.GetInstance().AddQuestionPack("алгебра")
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}

func TestAddQuestion(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	question := models.Question{
		Answers: []string{},
		Correct: 2,
		Text:    "",
		PackID:  1,
	}

	err = database.GetInstance().AddQuestion(question)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}

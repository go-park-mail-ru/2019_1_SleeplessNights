package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/jackc/pgx"
	"testing"
)

func TestAddMessageSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}

	uid := uint64(1)
	nickname := "test"
	avatarPath := "default_avatar.jpg"

	tId, err := database.GetInstance().UpdateTalker(uid, nickname, avatarPath)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	talkersArray := []uint64{1}

	rId, err := database.GetInstance().AddRoom(talkersArray)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	payload := `{"nickname":"Guest","avatarPath":"default_avatar.jpg","id":1,"text":"WOWOWOW"}`

	err = database.GetInstance().AddMessage(tId, rId, []byte(payload))
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}

func TestAddUserUnsuccessful_NoFoundTalker(t *testing.T) {

	talkersArray := []uint64{1}

	rId, err := database.GetInstance().AddRoom(talkersArray)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	payload := `{"nickname":"Guest","avatarPath":"default_avatar.jpg","id":1,"text":"WOWOWOW"}`

	err = database.GetInstance().AddMessage(999, rId, []byte(payload))
	if err == nil {
		t.Errorf("DB didn't return any error")
		return
	} else if err, _ := err.(pgx.PgError); err.Code != "23503" {
		t.Errorf("DB returned wrong error: %v",
			err.Error())
	}
}

func TestAddUserUnsuccessful_NoFoundRoom(t *testing.T) {

	uid := uint64(2)
	nickname := "test"
	avatarPath := "default_avatar.jpg"

	tId, err := database.GetInstance().UpdateTalker(uid, nickname, avatarPath)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	payload := `{"nickname":"Guest","avatarPath":"default_avatar.jpg","id":1,"text":"WOWOWOW"}`

	err = database.GetInstance().AddMessage(tId, 999, []byte(payload))
	if err == nil {
		t.Errorf("DB didn't return any error")
		return
	} else if err, _ := err.(pgx.PgError); err.Code != "23503" {
		t.Errorf("DB returned wrong error: %v",
			err.Error())
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}
}

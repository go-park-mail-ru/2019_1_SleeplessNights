package user_manager

import "testing"

func TestMakeSalt (t *testing.T){
	salt, err := MakeSalt()
	if err != nil{
		t.Errorf("MakeSalt returned error: %v", err.Error())
	}
	if len(salt) == 0{
		t.Errorf("Salt is empty")
	}
}

func TestMakePasswordHash(t *testing.T) {
	password := "1209qawsed"
	salt := []byte("qwehjsfksdf")
	hash := MakePasswordHash(password, salt)
	if len(hash) == 0{
		t.Errorf("Hash is empty")
	}
}
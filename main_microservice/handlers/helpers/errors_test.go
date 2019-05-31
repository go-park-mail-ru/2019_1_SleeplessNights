package helpers

import (
	"testing"
)

func TestMarshalToJSON(t *testing.T) {
	var errSet ErrorSet
	errSet = append(errSet, FormParsingErrorMsg)
	_, err := MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, UniqueEmailErrorMsg)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, MissedUserErrorMsg)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, InvalidEmailErrorMsg)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, WrongPassword)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, PasswordsDoNotMatchErrorMsg)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, PasswordIsTooSmallErrorMsg)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, InvalidNicknameErrorMsg)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, NicknameIsTooSmallErrorMsg)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, NicknameIsTooLongErrorMsg)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, AvatarExtensionError)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, AvatarIsMissingError)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	errSet = append(errSet, AvatarFileIsTooBig)
	_, err = MarshalToJSON(errSet)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
}


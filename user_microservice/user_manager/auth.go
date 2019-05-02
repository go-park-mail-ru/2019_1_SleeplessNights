package user_manager

import (
	"bytes"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

func (auth *authManager)CheckToken(ctx context.Context, in *services.SessionToken)(*services.User, error) {
	rawToken, err := jwt.Parse([]byte(in.Token))
	if err != nil {
		return nil, err
	}
	verifier := jwt.NewHMAC(jwt.SHA512, auth.secret)
	err = rawToken.Verify(verifier)
	if err != nil {
		return nil, err
	}
	payload := jwt.Payload{}
	_, err = rawToken.Decode(&payload)
	if err != nil {
		return nil, err
	}
	expValidator := jwt.ExpirationTimeValidator(time.Now(), true)
	err = payload.Validate(expValidator)
	if err != nil {
		return nil, err
	}
	userID, err := strconv.ParseUint(payload.JWTID, 10, 64)
	if err != nil {
		return nil, err
	}

	var user services.User
	user, err = database.GetInstance().GetUserByID(userID)
	if err !=nil {
		return nil, err
	}

	return &user, nil
}

func (auth *authManager)MakeToken(ctx context.Context, in *services.UserSignature)(*services.SessionToken, error) {
	id, password, salt, err := database.GetInstance().GetUserSignature(in.Email)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(password, MakePasswordHash(in.Password, salt)) {
		return nil, errors.AuthWrongPassword
	}

	signer := jwt.NewHMAC(jwt.SHA512, auth.secret)
	header := jwt.Header{}
	expiresAt := time.Now().Add(sessionLifeLen)
	payload := jwt.Payload{
		ExpirationTime: expiresAt.Unix(),
		JWTID:          strconv.FormatUint(id, 10),
	}
	token, err := jwt.Sign(header, payload, signer)
	if err != nil {
		return nil, err
	}

	var sessionToken services.SessionToken
	sessionToken.Token = string(token)
	return &sessionToken, nil
}
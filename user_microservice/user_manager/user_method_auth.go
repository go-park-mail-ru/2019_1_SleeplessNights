package user_manager

import (
	"bytes"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

const defaultSessionLifeLen  = 4 * time.Hour

func (us *userManager) CheckToken(ctx context.Context, in *services.SessionToken) (*services.User, error) {
	rawToken, err := jwt.Parse([]byte(in.Token))
	if err != nil {
		logger.Errorf("Failed to parse token: %v", err.Error())
		return nil, err
	}
	verifier := jwt.NewHMAC(jwt.SHA512, us.secret)
	err = rawToken.Verify(verifier)
	if err != nil {
		logger.Errorf("Failed to verify: %v", err.Error())
		return nil, err
	}
	payload := jwt.Payload{}
	_, err = rawToken.Decode(&payload)
	if err != nil {
		logger.Errorf("Failed to decode payload: %v", err.Error())
		return nil, err
	}
	expValidator := jwt.ExpirationTimeValidator(time.Now(), true)
	err = payload.Validate(expValidator)
	if err != nil {
		logger.Errorf("Failed to validate: %v", err.Error())
		return nil, err
	}
	userID, err := strconv.ParseUint(payload.JWTID, 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse payload: %v", err.Error())
		return nil, err
	}

	var user services.User
	user, err = database.GetInstance().GetUserByID(userID)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to get user: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}

	return &user, nil
}

func (us *userManager) MakeToken(ctx context.Context, in *services.UserSignature) (*services.SessionToken, error) {
	id, password, salt, err := database.GetInstance().GetUserSignature(in.Email)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to get signature: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}

	if !bytes.Equal(password, MakePasswordHash(in.Password, salt)) {
		logger.Errorf("Bytes isn't equal")
		return nil, errors.AuthWrongPassword
	}

	signer := jwt.NewHMAC(jwt.SHA512, us.secret)
	header := jwt.Header{}
	sessionLifeLen, err := time.ParseDuration(config.GetString("user_ms.pkg.user_manager.session_life_len"))
	if err != nil {
		sessionLifeLen = defaultSessionLifeLen
	}
	expiresAt := time.Now().Add(sessionLifeLen)
	payload := jwt.Payload{
		ExpirationTime: expiresAt.Unix(),
		JWTID:          strconv.FormatUint(id, 10),
	}
	token, err := jwt.Sign(header, payload, signer)
	if err != nil {
		logger.Errorf("Failed to sing: %v", err.Error())
		return nil, err
	}

	var sessionToken services.SessionToken
	sessionToken.Token = string(token)
	return &sessionToken, nil
}

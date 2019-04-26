package auth

import (
	"context"
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/services"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

const (
	sessionLifeLen = 4 * time.Hour
	NoTokenOwner = "error: There are no token's owner in database"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Auth")
	logger.SetLogLevel(logrus.TraceLevel)
}

var auth *authManager

type authManager struct {
	secret []byte
}

func init() {
	secretFile, err := os.Open(os.Getenv("BASEPATH") + "/secret")
	defer func() {
		err := secretFile.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	if err != nil {
		logger.Fatal(err)
	}

	var secret []byte
	_, err = fmt.Fscanln(secretFile, &secret)
	if err != nil {
		return
	}

	auth = &authManager{
		secret: secret,
	}
}

func GetInstance() *authManager {
	return auth
}

func (auth *authManager)Check(ctx context.Context, in *services.SessionToken)(*services.UserID, error) {
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

	var user services.UserID
	user.ID = userID

	return &user, nil
}

func (auth *authManager)MakeToken(ctx context.Context, in *services.UserID)(*services.SessionToken, error) {
	signer := jwt.NewHMAC(jwt.SHA512, auth.secret)
	header := jwt.Header{}
	expiresAt := time.Now().Add(sessionLifeLen)
	payload := jwt.Payload{
		ExpirationTime: expiresAt.Unix(),
		JWTID:          strconv.FormatUint(in.ID, 10),
	}
	token, err := jwt.Sign(header, payload, signer)
	if err != nil {
		return nil, err
	}

	var sessionToken services.SessionToken
	sessionToken.Token = string(token)
	return &sessionToken, nil
}

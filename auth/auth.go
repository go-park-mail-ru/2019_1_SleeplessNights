package auth

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	saltLen        = 16
	sessionLifeLen = 4 * time.Hour
	NoTokenOwner = "error: There are no token's owner in database"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Auth")
	logger.SetLogLevel(logrus.TraceLevel)
}

var secret []byte

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
	_, err = fmt.Fscanln(secretFile, &secret)
	if err != nil {
		return
	}
}

func MakeSalt() (salt []byte, err error) {
	salt = make([]byte, saltLen)
	_, err = rand.Read(salt) //Заполняем слайс случайными значениями по всей его длине
	return
}

func MakePasswordHash(password string, salt []byte) (hash []byte) {
	saltedPassword := bytes.Join([][]byte{[]byte(password), salt}, nil)
	hashedPassword := sha512.Sum512(saltedPassword) //sha512 возвращает массив, а слайс можно взять только по addressable массиву
	hash = hashedPassword[0:]
	return
}

func MakeSession(user models.User) (sessionCookie http.Cookie, err error) {
	signer := jwt.NewHMAC(jwt.SHA512, secret)
	header := jwt.Header{}
	expiresAt := time.Now().Add(sessionLifeLen)
	payload := jwt.Payload{
		ExpirationTime: expiresAt.Unix(),
		JWTID:          strconv.FormatUint(uint64(user.ID), 10),
	}
	token, err := jwt.Sign(header, payload, signer)
	if err != nil {
		return
	}

	//Создаём cookie сессии
	sessionCookie = http.Cookie{
		Name:     "session_token",
		Value:    string(token),
		Expires:  expiresAt,
		HttpOnly: true,
	}
	return
}

func Authorize(sessionToken string) (user models.User, err error) {
	fmt.Println("HERE")
	logger.Debug("Authorize method called with token:",sessionToken)
	rawToken, err := jwt.Parse([]byte(sessionToken))
	if err != nil {
		return
	}
	verifier := jwt.NewHMAC(jwt.SHA512, secret)
	err = rawToken.Verify(verifier)
	if err != nil {
		logger.Debug("Verify function returned an error:",err)
		return
	}
	payload := jwt.Payload{}
	_, err = rawToken.Decode(&payload)
	if err != nil {
		logger.Debug("Decode function returned an error:",err)
		return
	}
	expValidator := jwt.ExpirationTimeValidator(time.Now(), true)
	err = payload.Validate(expValidator)
	if err != nil {
		logger.Debug("Validate function returned an error:",err)
		return
	}
	userID, err := strconv.ParseUint(payload.JWTID, 10, 64)
	if err != nil {
		logger.Debug("ParseUint function returned an error:",err)
		return
	}
	user, err = database.GetInstance().GetUserViaID(uint(userID))
	if err != nil {
		logger.Debug("GetUserViaID returned an error:",err)
		if err.Error() == database.NoUserFound {
			return user, errors.New(NoTokenOwner)
		} else {
			return
		}
	}
	return
}

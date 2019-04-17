package auth

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"log"
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

var secret []byte

func init() {

	secretFile, err := os.Open(os.Getenv("BASEPATH") + "/secret")
	defer func() {
		err := secretFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
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
	rawToken, err := jwt.Parse([]byte(sessionToken))
	if err != nil {
		return
	}
	verifier := jwt.NewHMAC(jwt.SHA512, secret)
	err = rawToken.Verify(verifier)
	if err != nil {
		return
	}
	payload := jwt.Payload{}
	_, err = rawToken.Decode(&payload)
	if err != nil {
		return
	}
	expValidator := jwt.ExpirationTimeValidator(time.Now(), true)
	err = payload.Validate(expValidator)
	if err != nil {
		return
	}
	userID, err := strconv.ParseUint(payload.JWTID, 10, 64)
	if err != nil {
		return
	}
	user, err = database.GetInstance().GetUserViaID(uint(userID))
	if err != nil {
		if err.Error() == database.NoUserFound {
			return user, errors.New(NoTokenOwner)
		} else {
			return
		}
	}
	return
}

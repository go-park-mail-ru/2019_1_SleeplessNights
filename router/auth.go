package router

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var secret []byte

func init() {
	secretFile, err := os.Open("secret")
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fscanln(secretFile, &secret)
	if err != nil {
		return
	}
}

func MakeSalt()(salt []byte, err error) {
	saltLen, err := strconv.Atoi(os.Getenv("SaltLen"))
	if err != nil {
		return
	}
	salt = make([]byte, saltLen)
	rand.Read(salt)//Заполняем слайс случайными значениями по всей его длине
	return
}

func MakePasswordHash(password string, salt []byte)(hash []byte){
	saltedPassword:= bytes.Join([][]byte {[]byte(password), salt}, nil)
	hashedPassword := sha512.Sum512(saltedPassword)//sha512 возвращает массив, а слайс можно взять только по addressable массиву
	hash = hashedPassword[0:]
	return
}

func MakeSession(user *models.User)(sessionCookie http.Cookie, err error){
	signer := jwt.NewHMAC(jwt.SHA512, secret)
	header := jwt.Header{}
	sessionLifeLen, err := time.ParseDuration(os.Getenv("SessionLifeLen"))
	if err != nil {
		return
	}
	expiresAt := time.Now().Add(sessionLifeLen)
	payload := jwt.Payload{
		ExpirationTime: expiresAt.Unix(),
		JWTID: strconv.FormatUint(uint64(user.ID), 10),
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

func Authorize(sessionToken string)(*models.User, error){
	rawToken, err := jwt.Parse([]byte(sessionToken))
	if err != nil {
		fmt.Println("Error while parsing token")
		return nil, err
	}
	verifier := jwt.NewHMAC(jwt.SHA512, secret)
	err = rawToken.Verify(verifier)
	if err != nil {
		fmt.Println("Error while verifying token")
		return nil, err
	}
	payload := jwt.Payload{}
	_, err = rawToken.Decode(&payload)
	if err != nil {
		fmt.Println("Error while decoding token")
		fmt.Println(err)
		return nil, err
	}
	expValidator := jwt.ExpirationTimeValidator(time.Now(), true)
	err = payload.Validate(expValidator)
	if err != nil {
		fmt.Println("Error while validating token")
		return nil, err
	}
	userID, err := strconv.ParseUint(payload.JWTID, 10, 64)
	if err != nil {
		return nil, err
	}
	user, found := models.Users[models.UserKeyPairs[uint(userID)]]
	if !found {
		fmt.Println("Error: user not found in the database")
		return nil, errors.New("error: There are no token's owner in database")
	}
	return &user, nil
}
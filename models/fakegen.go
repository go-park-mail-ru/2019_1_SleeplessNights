package models

import (
	"bytes"
	"crypto/sha512"
	"github.com/manveru/faker"
	"math/rand"
	"time"
)

const (
	saltLen        = 16
	sessionLifeLen = 4 * time.Hour
)

func MakeSalt() (salt []byte, err error) {
	salt = make([]byte, saltLen)
	rand.Read(salt) //Заполняем слайс случайными значениями по всей его длине
	return
}

func MakePasswordHash(password string, salt []byte) (hash []byte) {
	saltedPassword := bytes.Join([][]byte{[]byte(password), salt}, nil)
	hashedPassword := sha512.Sum512(saltedPassword) //sha512 возвращает массив, а слайс можно взять только по addressable массиву
	hash = hashedPassword[0:]
	return
}

// Fills Users Map with userdata
func CreateFakeData(Quantity int) {
	Fake, _ := faker.New("en")

	Salt, _ := MakeSalt()
	Password := "1Q2W3e4r5t6y7u"
	TimeDuration, _ := time.ParseDuration("1h")

	for i := 1; i <= Quantity; i++ {
		user := User{
			ID:         uint(i),
			Email:      Fake.Email(),
			Password:   MakePasswordHash(Password, Salt),
			Salt:       Salt,
			Won:        uint(rand.Uint32()),
			Lost:       uint(rand.Uint32()),
			PlayTime:   TimeDuration,
			Nickname:   Fake.UserName(),
			AvatarPath: "default_avatar.jpg",
		}
		Users[user.Email] = user
		UserKeyPairs[user.ID] = user.Email
	}

}

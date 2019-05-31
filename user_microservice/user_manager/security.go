package user_manager

import (
	"bytes"
	"crypto/sha512"
	"math/rand"
)

const saltLen = 16

func MakeSalt() (salt []byte, err error) {
	salt = make([]byte, saltLen)
	_, err = rand.Read(salt) //Заполняем слайс случайными значениями по всей его длине
	if err != nil {
		logger.Errorf("Failed to read salt: %v", err.Error())
		return
	}
	return
}

func MakePasswordHash(password string, salt []byte) (hash []byte) {
	saltedPassword := bytes.Join([][]byte{[]byte(password), salt}, nil)
	hashedPassword := sha512.Sum512(saltedPassword) //sha512 возвращает массив, а слайс можно взять только по addressable массиву
	hash = hashedPassword[0:]
	return
}

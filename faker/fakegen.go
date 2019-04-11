package faker

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/manveru/faker"
	"math/rand"
)

const (
	FakeUserPassword = "1Q2W3e4r5t6y7u"
)

// Fills Users Map with user data
func CreateFakeData(quantity int) {
	fake, err := faker.New("en")
	if err != nil {
		logger.Fatal.Println(err.Error())
		return
	}
	for i := 1; i <= quantity; i++ {
		salt, err := helpers.MakeSalt()
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}

		email := fake.Email()
		_, err = database.GetInstance().GetUserViaEmail(email)
		if err == nil {
			continue
		}
		user := models.User{
			ID:       	models.MakeID(),
			Email:    	email,
			Password: 	helpers.MakePasswordHash(FakeUserPassword, salt),
			Salt:       salt,
			Won:        uint(rand.Uint32()),
			Lost:       uint(rand.Uint32()),
			PlayTime:   0,
			Nickname:   fake.UserName(),
			AvatarPath: "default_avatar.jpg",
		}
		_ = database.GetInstance().AddUser(user)
	}
}
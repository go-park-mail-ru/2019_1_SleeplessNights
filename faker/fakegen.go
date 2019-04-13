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
	FakeUserPassword             = "1Q2W3e4r5t6y7u"
	NumberOfPacks                = 10
	NumberOfQuestionsInOnePack   = 10
	NumberOfAnswersInOneQuestion = 4
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
			Email:      email,
			Password:   helpers.MakePasswordHash(FakeUserPassword, salt),
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

func CreateFakePacks() {
	fake, err := faker.New("en")
	if err != nil {
		logger.Fatal.Println(err.Error())
		return
	}
	themes := []string{
		"математика",
		"информатика",
		"химия",
		"биология",
		"физика",
		"культура",
		"история",
		"языки",
		"философия",
		"мемология",
	}

	var packID uint = 1
	for i := 0; i < NumberOfPacks; i++ {
		for _, theme := range themes {
			for i := 0; i < NumberOfQuestionsInOnePack; i++ {
				var answers []string
				for j := 0; j < NumberOfAnswersInOneQuestion; j++ {
					answer := fake.CompanyName()
					answers = append(answers, answer)
				}
				question := models.Question{
					Answers: answers,
					Correct: 1,
					Text:    fake.Paragraph(1, true),
					PackID:  packID,
				}

				err := database.GetInstance().AddQuestion(question)
				if err != nil {
					logger.Error.Println(err.Error())
					return
				}
			}
			err := database.GetInstance().AddQuestionPack(theme)
			if err != nil {
				logger.Error.Println(err.Error())
				return
			}
			packID++
		}
	}
}

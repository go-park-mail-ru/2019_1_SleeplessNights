package faker

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database/models"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/manveru/faker"
	"github.com/sirupsen/logrus"
)

const (
	FakeUserPassword             = "1Q2W3e4r5t6y7u"
	NumberOfPacks                = 10
	NumberOfQuestionsInOnePack   = 10
	NumberOfAnswersInOneQuestion = 4
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("GameMS_Faker")
	logger.SetLogLevel(logrus.TraceLevel)
}

func CreateFakePacks() {
	fake, err := faker.New("en")
	if err != nil {
		logger.Error(err.Error())
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
				logger.Error(err.Error())
				return
			}
		}
		err := database.GetInstance().AddQuestionPack(theme)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		packID++
	}
}

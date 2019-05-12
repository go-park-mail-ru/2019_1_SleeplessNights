package database

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database/models"
	"os"
)

func loadQuestionsJson(file string) (questions []models.Question, ) {
	questionJson, err := os.Open(file)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	jsonParser := json.NewDecoder(questionJson)
	err = jsonParser.Decode(&questions)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = questionJson.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

func loadPacksJson(file string) (packs []models.Pack, ) {
	questionJson, err := os.Open(file)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	jsonParser := json.NewDecoder(questionJson)
	err = jsonParser.Decode(&packs)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = questionJson.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

func (db *dbManager) PopulateDatabase() {
	packs := loadPacksJson(os.Getenv("BASEPATH") + "/game_microservice/database/packs.json")
	for _, pack := range packs {
		err := db.AddQuestionPack(pack.Theme)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}

	questions := loadQuestionsJson(os.Getenv("BASEPATH") + "/game_microservice/database/questions.json")
	for _, question := range questions {
		err := db.AddQuestion(question)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
}

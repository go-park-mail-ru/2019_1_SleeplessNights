package database

import (
	"encoding/json"
	"os"
)

func loadQuestionsJson(file string) (questions []Question, ) {
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

func loadPacksJson(file string) (packs []Pack, ) {
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
		err := db.AddQuestionPack(pack)
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

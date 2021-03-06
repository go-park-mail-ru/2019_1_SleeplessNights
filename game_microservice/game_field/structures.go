package game_field

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"

type gameCell struct {
	isAvailable  bool
	answerResult int
	question     *database.Question
}

type Pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type gfPlayer struct {
	pos          *Pair //Поставил указатель на Pair, чтобы pos поддерживала значение nil (начальные условия)
	rightAnswers int
	falseAnswers int
	partyCounter int
	id           uint64
}

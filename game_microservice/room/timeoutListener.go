package room

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
)

func (r *Room) AnswerTimerFunc() {

	logger.Info("Timer Has Expired before answer was given")

	result, err := json.Marshal(message.Message{message.ClientAnswer, message.Answer{-1}})
	if err != nil {
		logger.Error("startAnswerTimeChecker Error trying to marshall Answer response")
	}
	var answer message.Message
	err = json.Unmarshal(result, &answer)
	if err != nil {
		logger.Error("startAnswerTimeChecker Error trying to UnMarshal Answer response")
	}
	r.requestsQueue <- MessageWrapper{r.active, answer}
}

func (r *Room) GoToTimerFunc() {

	logger.Info("Timer Has Expired before answer was given")

	result, err := json.Marshal(message.Message{message.GoTo, message.Coordinates{-1, -1}})
	if err != nil {
		logger.Error("startGoToTimeChecker Error trying to marshall GoTo response")
	}
	var GoToMessage message.Message
	err = json.Unmarshal(result, &GoToMessage)
	if err != nil {
		logger.Error("startGoToTimeChecker Error trying to UnMarshal GoTo response")
	}
	r.requestsQueue <- MessageWrapper{r.active, GoToMessage}
}

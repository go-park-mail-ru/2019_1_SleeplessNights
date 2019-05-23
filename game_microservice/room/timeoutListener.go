package room

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
)

func (r *Room) AnswerTimerFunc() {
	logger.Info("AnswerTimerFunc Timer Has Expired before answer was given")

	result, err := json.Marshal(message.Message{
		Title:   message.ClientAnswer,
		Payload: message.Answer{AnswerId: -1},
	})
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

	logger.Info("GoToTimerFunc Timer Has Expired before answer was given")

	result, err := json.Marshal(message.Message{
		Title:   message.GoTo,
		Payload: message.Coordinates{X: -1, Y: -1},
	})
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

func (r *Room) ChoosePackTimerFunc() {

	logger.Info("ChoosePackTimerFunc Timer Has Expired before answer was given")

	result, err := json.Marshal(message.Message{
		Title: message.NotDesiredPack,
		Payload: message.PackID{PackId: -1},
	})
	if err != nil {
		logger.Error("ChoosePackTimerFunc Error trying to marshall NotDesiredPack response")
	}
	var PackIdMessage message.Message
	err = json.Unmarshal(result, &PackIdMessage)
	if err != nil {
		logger.Error("ChoosePackTimerFunc Error trying to UnMarshal NotDesiredPack response")
	}
	r.requestsQueue <- MessageWrapper{r.active, PackIdMessage}
}

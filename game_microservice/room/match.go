package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game_field"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/messge"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"time"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func (r *Room) buildEnv() {
	packs, err := database.GetInstance().GetPacksOfQuestions(10)
	if err != nil {
		logger.Error("Error occurred while fetching question packs from DB:", err)
		//TODO deal with error, maybe kill the room
	}
	packIDs := make([]int, len(packs))
	for _, pack := range packs {
		packIDs = append(packIDs, int(pack.ID))
	}

	questions, err := database.GetInstance().GetQuestions(packIDs)
	if err != nil || len(questions) < game_field.QuestionsNum {
		logger.Error("Error occurred while fetching question from DB:", err)
		//TODO deal with error, maybe kill the room
	}
	/*var localQuestions [game_field.QuestionsNum]local.Question
	var lq local.Question
	for i := 0; i < len(localQuestions); i++ {
		questionJSON, err := json.Marshal(questions[i])
		if err != nil {
			logger.Error("Error occurred while marshalling question into JSON:", err)
			//TODO deal with error, maybe refresh questions
		}
		lq = local.Question{PackID: uint64(questions[i].ID),
			QuestionJson:    string(questionJSON),
			CorrectAnswerId: questions[i].Correct}
		localQuestions[i] = lq
	}
	*/
	r.field.Build(questions)
	//Процедура должна пересоздавать игровое поле, запрашивать новый список тем из БД и готовить комнату к новой игре
	//При этом она должна уметь работать асинхронно и не выбрасывать пользователей из комнаты во время работы
}

// TODO PREPAREMATCH AND BUILD ENV (simultaneously (optional), then wait them both to work out, use with WaitGroup )

func (r *Room) prepareMatch() {
	logger.Info("Entered Prepare Match")
	r.buildEnv()
	r.requestsQueue = make(chan MessageWrapper, channelCapacity)
	r.responsesQueue = make(chan MessageWrapper, channelCapacity)

	p1Chan := r.p1.Subscribe()
	p2Chan := r.p2.Subscribe()

	err := r.notifyAll(messge.Message{Title: messge.StartGame, Payload: nil})
	if err != nil {
		logger.Error("Failed to notify all players:", err)
	}
	logger.Info("Игрокам Отправлены StartGame")
	r.waitForSyncMsg = messge.Ready
	//Read Messages from Players
	//Moved message receive conditions to Requests handler

	go func() {
		for msgP1 := range p1Chan {
			logger.Info("got message from P1", msgP1)
			r.requestsQueue <- MessageWrapper{&r.p1, msgP1}
		}
	}()

	go func() {
		for msgP2 := range p2Chan {
			logger.Info("got message from P2", msgP2)
			r.requestsQueue <- MessageWrapper{&r.p2, msgP2}
		}
	}()
	//Channel to Write Server messages to the player1/player2
	r.startMatch()
}

func (r *Room) startMatch() {
	//   Эта процедура запускает игровой процесс
	//   Здесь мы будем слушать все сообщения пользователей асинхронно и складывать их в очередь для обработки
	//   В цикле мы будем обрабатывать все входные сообщения, выполнять нашу бизнес логику (менять значение таймера,
	//   давать пользователям вопросы и т.д.)
	//   Все сообщения мы будем складывать в очередь на отправку и отправлять всю очередь каждые 0.5 сек
	//   (цифра примерная, может поменяться и должна быть вынесена в костанту)
	//TODO develop

	// Call Prepare Room

	logger.Infof("StartMatch : Game process has started p1: %d, p2: %d", r.p1.UID(), r.p2.UID())

	go func() {
		for serverResponse := range r.responsesQueue {
			logger.Info("Got message to Send", serverResponse)
			err := (*serverResponse.player).Send(serverResponse.msg)
			if err != nil {
				logger.Error("responseQueue: error trying to send response to player", err)
			}
		}
		time.Sleep(responseInterval * time.Millisecond)
	}()

	go func() {
		for msg := range r.requestsQueue {
			//Проверка структуры сообщения
			logger.Info("Got Message from client")
			if !msg.msg.IsValid() {
				logger.Error("Message to send has invalid structure")
				continue
			}

			if !r.isSyncValid(msg) {
				logger.Warningf("Got message of type %s from player %d, expected %s from player %d",
					msg.msg.Title, msg.player, r.waitForSyncMsg, r.active)
				continue
			}
			logger.Info("Message entered mux")
			r.MessageHandlerMux(msg)
		}
	}()
}

package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game_field"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

var logger *log.Logger

func init() {

	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

var userManager services.UserMSClient

func init() {
	var err error
	grpcConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Fatal("Can't connect to auth microservice via grpc")
	}
	userManager = services.NewUserMSClient(grpcConn)
	closer.Bind(func() {
		err := grpcConn.Close()
		if err != nil {
			logger.Error("Error occurred while closing grpc connection", err)
		}
	})
}

//Build Environment After getting desiredPacks
func (r *Room) buildEnv() {
	logger.Info("Entered BuildEnv in Room")

	packs, err := database.GetInstance().GetPacksOfQuestions(6)

	if err != nil {
		logger.Error("Error occurred while fetching question packs from DB:", err)
		//TODO deal with error, maybe kill the room
	}
	logger.Info("Got packs from database")
	packIDs := make([]uint64, 0)
	for _, pack := range packs {
		packIDs = append(packIDs, uint64(pack.ID))
		fieldPacks := r.field.GetThemesSlice()
		*fieldPacks = append(*fieldPacks, message.ThemePack{pack.ID, pack.Theme})
	}

	questions, err := database.GetInstance().GetQuestions(packIDs)
	if err != nil || len(questions) < game_field.QuestionsNum {
		logger.Error("Error occurred while fetching question from DB:", err)
		//TODO deal with error, maybe kill the room
	}

	r.field.Build(questions)
	//Процедура должна пересоздавать игровое поле, запрашивать новый список тем из БД и готовить комнату к новой игре
	//При этом она должна уметь работать асинхронно и не выбрасывать пользователей из комнаты во время работы
}

// TODO PREPAREMATCH AND BUILD ENV (simultaneously (optional), then wait them both to work out, use with WaitGroup )

func (r *Room) prepareMatch() {

	logger.Info("Entered Prepare Match Room")
	r.buildEnv()

	r.requestsQueue = make(chan MessageWrapper, channelCapacity)
	r.responsesQueue = make(chan MessageWrapper, channelCapacity)

	p1Chan := r.p1.Subscribe()
	p2Chan := r.p2.Subscribe()

	err := r.notifyAll(message.Message{Title: message.StartGame, Payload: nil})
	if err != nil {
		logger.Error("Failed to notify about StartGame ,to all players:", err)
	}

	user2, err := userManager.GetUserById(context.Background(), &services.UserId{Id: r.p2.UID()})
	if err != nil {
		logger.Error("failed to get userprofile2 from grpc:", err)
	}
	err = r.notifyP1(message.Message{Title: message.OpponentProfile, Payload: user2})
	if err != nil {
		logger.Error("Failed to notify Player 1:", err)
	}
	user1, err := userManager.GetUserById(context.Background(), &services.UserId{Id: r.p2.UID()})
	if err != nil {
		logger.Error("failed to get userprofile1 from grpc:", err)
	}
	err = r.notifyP2(message.Message{Title: message.OpponentProfile, Payload: user1})
	if err != nil {
		logger.Error("Failed to notify Player 2:", err)
	}

	logger.Info("Игрокам Отправлены StartGame")

	err = r.notifyAll(message.Message{Title: message.Themes, Payload: r.field.GetThemesSlice()})
	if err != nil {
		logger.Error("Failed to send ThemesResponse , to all players:", err)
	}
	packArray := r.field.GetQuestionsThemes()
	err = r.notifyAll(message.Message{Title: message.QuestionsThemes, Payload: packArray})

	if err != nil {
		logger.Error("Failed to send QuestionsThemesResponse , to all players:", err)
	}
	r.waitForSyncMsg = message.Ready
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

//Вызов после отправления READY игрокам
func (r *Room) startMatch() {
	//   Эта процедура запускает игровой процесс
	//   Здесь мы будем слушать все сообщения пользлователей асинхронно и складывать их в очередь для обработки
	//   В цикле мы будем обрабатывать все входные сообщения, выполнять нашу бизнес логику (менять значение таймера,
	//   давать пользователям вопросы и т.д.)
	//   Все сообщения мы будем складывать в очередь на отправку и отправлять всю очередь каждые 0.5 сек
	//   (цифра примерная, может поменяться и должна быть вынесена в костанту)
	//TODO develop

	// Call Prepare Room

	logger.Infof("StartMatch : Game process has started p1 UID: %d, p2 UID: %d", r.p1.UID(), r.p2.UID())

	go func() {
		for serverResponse := range r.responsesQueue {
			logger.Info("Got message to Send recepient: UID", (*serverResponse.player).UID(), "Message:", serverResponse.msg)
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
				logger.Error("Got message with invalid structure")
				continue
			}

			if !r.isSyncValid(msg) {
				logger.Warningf("Got SyncInvalid message of type %s from player UID %d, expected %s from player %d",
					msg.msg.Title, (*msg.player).UID(), r.waitForSyncMsg, r.active)
				continue
			}
			logger.Info("Message entered mux")
			r.MessageHandlerMux(msg)
		}
	}()
}

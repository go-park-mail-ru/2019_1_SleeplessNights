package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game_field"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player"
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

	logger = log.GetLogger("GameMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

var userManager services.UserMSClient
const (
	StartGameDelay = 1100
)
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

	packs := r.field.GetPacksSlice()
	logger.Info("Got packs from database")
	packIDs := make([]uint64, 0)
	for _, pack := range *packs {
		packIDs = append(packIDs, uint64(pack.ID))
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
	
	//Где-то здесь добавить выбор паков игроками

	logger.Info("Entered Prepare Match Room")
	logger.Info("Delay")
	time.Sleep(StartGameDelay*time.Microsecond)
	//BuildEnv достает только выбранные паки и строит игровое поле по ним
	r.buildEnv()

	//Сюда приходим после тогос как паки будут выбраны
	err := r.notifyAll(message.Message{Title: message.StartGame, Payload: nil})
	if err != nil {
		logger.Error("Failed to notify about StartGame ,to all players:", err)
	}

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

	//err = r.notifyAll(message.Message{Title: message.Themes, Payload: r.field.GetThemesSlice()})
	//if err != nil {
	//	logger.Error("Failed to send ThemesResponse , to all players:", err)
	//}
	packArray := r.field.GetQuestionsThemes()
	err = r.notifyAll(message.Message{Title: message.QuestionsThemes, Payload: packArray})

	if err != nil {
		logger.Error("Failed to send QuestionsThemesResponse , to all players:", err)
	}
	//Read Messages from Players
	//Moved message receive conditions to Requests handler
	r.waitForSyncMsg = message.GoTo
	r.responsesQueue <- MessageWrapper{&r.p1, message.Message{Title: message.YourTurn, Payload: nil}}
	r.responsesQueue <- MessageWrapper{&r.p2, message.Message{Title: message.OpponentTurn, Payload: nil}}
	// Результат работы достаем из канала Events()отсылаем в канал ResponsesQueue
	cellsSlice := r.field.GetAvailableCells(r.getPlayerIdx(r.active))

	var secondPlayer *player.Player
	if &r.p1 == r.active {
		secondPlayer = &r.p2
	}
	if &r.p2 == r.active {
		secondPlayer = &r.p1
	}
	cells := make([]Pair, 0)
	for _, cell := range cellsSlice {
		cells = append(cells, Pair{cell.X, cell.Y})
	}
	payload := struct {
		CellsSlice []Pair
		Time       time.Duration
	}{
		CellsSlice: cells,
		Time:       timeToMove,
	}
	//Send Available cells to active player (Do it every time, after giving player a turn rights
	r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.AvailableCells, Payload: payload}}
	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.AvailableCells, Payload: payload}}
	r.timerToMove = time.AfterFunc(timeToMove*time.Second, r.GoToTimerFunc)
}

//Точка входа в игровой процесс
func (r *Room) startGameProcess() {

	user2, err := userManager.GetUserById(context.Background(), &services.UserId{Id: r.p2.UID()})
	if err != nil {
		logger.Error("failed to get userprofile2 from grpc:", err)
	}
	err = r.notifyP1(message.Message{Title: message.OpponentProfile, Payload: user2})

	//Send available pack to players
	packs, err := database.GetInstance().GetPacksOfQuestions(packTotal)
	logger.Info("packs", packs)
	if err != nil {
		logger.Error("Failed to get available packs from database")
	}
	fieldPacks := r.field.GetPacksSlice()
	*fieldPacks = make([]database.Pack, packTotal)
	copy(*fieldPacks, packs)
	logger.Info(&r.p1, "  ", &r.p2)
	logger.Info("Packs", packs)
	payload := struct {
		Packs []database.Pack
		Time  time.Duration
	}{
		Packs: packs,
		Time:  timeToChoosePack,
	}
	r.responsesQueue <- MessageWrapper{&r.p1, message.Message{Title: message.AvailablePacks, Payload: payload}}
	r.responsesQueue <- MessageWrapper{&r.p2, message.Message{Title: message.AvailablePacks, Payload: payload}}

	r.responsesQueue <- MessageWrapper{&r.p1, message.Message{Title: message.YourTurn, Payload: nil}}
	r.responsesQueue <- MessageWrapper{&r.p2, message.Message{Title: message.OpponentTurn, Payload: nil}}
	r.waitForSyncMsg = message.NotDesiredPack
	r.timerToChoosePack = time.AfterFunc(timeToChoosePack*time.Second, r.ChoosePackTimerFunc)
	logger.Infof("StartMatch : Game process has started p1 UID: %d, p2 UID: %d", r.p1.UID(), r.p2.UID())
}

func (r *Room) changeTurn() {
	if (*r.active).ID() == r.p1.ID() {
		r.active = &r.p2
	} else {
		r.active = &r.p1
	}
}

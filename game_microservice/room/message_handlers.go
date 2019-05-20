package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/event"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player"
	"time"
)

type Pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (r *Room) ReadyHandler(m MessageWrapper) bool {

	if m.player == &r.p1 {
		logger.Info("Игрок 1 Готов")
		r.p1Status = StatusReady
	}
	if m.player == &r.p2 {
		logger.Info("Игрок 2 готов")
		r.p2Status = StatusReady
	}

	if r.p1Status == StatusReady && r.p2Status == StatusReady {
		//After getting ready messages from both players set p1 as active and send messages
		r.waitForSyncMsg = message.GoTo
		r.active = &r.p1

		logger.Info("Ход игрока 1, ожидание команды GoTo")

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

	return true
}

func (r *Room) GoToHandler(m MessageWrapper) bool {
	logger.Infof("player UID %d requested GoTo", (*m.player).UID())

	if r.timerToMove.Stop() {
		logger.Info("GoToHandler Timer is disabled manually")
	} else {
		logger.Info("GoToHandler Timer is disabled by timeout")
	}
	var secondPlayer *player.Player

	if &r.p1 == m.player {
		secondPlayer = &r.p2
	}
	if &r.p2 == m.player {
		secondPlayer = &r.p1
	}

	st := m.msg.Payload.(map[string]interface{})
	nextX := int(st["x"].(float64))
	nextY := int(st["y"].(float64))
	logger.Info("Sending SelectedCell Index to players")

	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.SelectedCell, Payload: message.Coordinates{nextX, nextY}}}
	r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.SelectedCell, Payload: message.Coordinates{nextX, nextY}}}

	//Признак таймаута хода игрока
	if nextY == -1 && nextX == -1 {
		//Смена хода
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.OpponentTurn, Payload: nil}}
		r.changeTurn()
		r.waitForSyncMsg = message.GoTo
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.YourTurn, Payload: nil}}

		cellsSlice := r.field.GetAvailableCells(r.getPlayerIdx(r.active))

		if &r.p1 == r.active {
			secondPlayer = &r.p2
		}
		if &r.p2 == r.active {
			secondPlayer = &r.p1
		}
		if len(cellsSlice) != 0 {

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
			r.timerToMove = time.AfterFunc(timeToAnswer*time.Second, r.GoToTimerFunc)
			r.waitForSyncMsg = message.GoTo
			return true
		} else {
			logger.Error("Unexpected condition")
		}
	}

	var eventSlice []event.Event
	var err error
	if &r.p1 == m.player {
		eventSlice, err = r.field.TryMovePlayer1(m.msg)
		secondPlayer = &r.p2
	}
	if &r.p2 == m.player {
		eventSlice, err = r.field.TryMovePlayer2(m.msg)
		secondPlayer = &r.p1
	}

	if err != nil {
		logger.Error("GoToHandler, called TryMovePLayer, got error", err)

		return false
	}

	for _, e := range eventSlice {
		if e.Etype == event.Info {
			logger.Info("player", (*r.active).ID(), "got question", e.Edata)
			q, ok := e.Edata.(string)
			if !ok {
				logger.Error("Go_To handler couldn't cast Edata interface with question to string")

				return false
			}
			r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.YourQuestion, Payload: q}}
			r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.OpponentQuestion, Payload: q}}
			r.waitForSyncMsg = "ANSWER"
			r.timerToAnswer = time.AfterFunc(timeToAnswer*time.Second, r.AnswerTimerFunc)
		}
		if e.Etype == event.WinPrize {
			//Write to DB results of the
			logger.Info("player", (*r.active).ID(), "Has Won the prize")
			r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.Win, Payload: nil}}
			r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.Loss, Payload: nil}}
			r.waitForSyncMsg = "Leave"
			r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.WannaPlayAgain, Payload: nil}}
			r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.WannaPlayAgain, Payload: nil}}
		}
	}

	return true
}

func (r *Room) ClientAnswerHandler(m MessageWrapper) bool {

	logger.Infof("player UID %d answered to ClientAnswerHandler", (*m.player).UID())

	if r.timerToAnswer.Stop() {
		logger.Info("ClientAnswerHandler Timer is disabled manually")
	} else {
		logger.Info("ClientAnswerHandler Timer is disabled by timeout")
	}
	logger.Info("")
	st, ok := m.msg.Payload.(map[string]interface{})
	if !ok {
		logger.Error("ClientAnswerHandler, couldn't cast payload with answer_id to map[string]interface{}")
	}
	answerId, ok := st["answer_id"].(float64)
	if !ok {
		logger.Error(`ClientAnswerHandler, couldn't find value in map st with key "answer_id" `)
	}
	q := r.field.GetRegisterQuestion()

	playerHasNoMoves := false

	if !r.field.CheckAnswer(int(answerId)) {
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.YourAnswer, Payload: message.AnswerResult{int(answerId), q.Correct}}}
		playerIdx := r.getPlayerIdx(r.active)
		if !r.field.CheckIfMovesAvailable(playerIdx) {
			playerHasNoMoves = true
		}
	} else {
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.YourAnswer, Payload: message.AnswerResult{int(answerId), q.Correct}}}
		r.field.Move(r.getPlayerIdx(r.active))
	}
	var secondPlayer *player.Player

	if &r.p1 == r.active {
		secondPlayer = &r.p2
	}
	if &r.p2 == r.active {
		secondPlayer = &r.p1
	}

	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.OpponentAnswer, Payload: message.AnswerResult{int(answerId), q.Correct}}}

	if playerHasNoMoves {
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.Loss, Payload: nil}}
		r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.Win, Payload: nil}}
		r.waitForSyncMsg = "Leave"
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.WannaPlayAgain, Payload: nil}}
		r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.WannaPlayAgain, Payload: nil}}

	} else {
		//Смена хода после ответа игрока
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.OpponentTurn, Payload: nil}}
		r.changeTurn()
		r.waitForSyncMsg = message.GoTo
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.YourTurn, Payload: nil}}

		cellsSlice := r.field.GetAvailableCells(r.getPlayerIdx(r.active))

		if &r.p1 == r.active {
			secondPlayer = &r.p2
		}
		if &r.p2 == r.active {
			secondPlayer = &r.p1
		}
		if len(cellsSlice) != 0 {

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

		} else {
			logger.Error("Unexpected condition")
		}
	}
	return true
}

func (r *Room) LeaveHandler(m MessageWrapper) bool {

	return true
}

//Оставить комнату с теми же игроками, создать для них новое игровое поле
//Если один из них голосует выйти, то написать об этом другому
func (r *Room) ContinueHandler(m MessageWrapper) bool {

	var secondPlayer *player.Player

	if &r.p1 == m.player {
		r.p1Status = StatusWannaContinue
		secondPlayer = &r.p2
	}
	if &r.p2 == m.player {
		r.p2Status = StatusWannaContinue
		secondPlayer = &r.p1
	}
	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.OpponentContinues, Payload: nil}}
	//Если оба игрока согласны продолжить игру, то при получении последнего
	// "WannaContinue" собираем игровое поле заново с другими вопросами

	if r.p2Status == StatusWannaContinue && r.p1Status == StatusWannaContinue {
		logger.Info("Both players wanna continue gaming, preparing new Match")
		r.field.ResetPlayersPositions()
		// Получить новый пак вопросов, заново заполнить ячейками игровое поле
		// Поставить состояние  игрового цикла в начало
		logger.Info("Building new environment")
		r.buildEnv()
		r.p1Status = StatusReady
		r.p2Status = StatusReady
		//Здесь перезупускаем игровой процесс с теми же игроками
		r.responsesQueue <- MessageWrapper{&r.p1, message.Message{Title: message.OpponentTurn, Payload: nil}}
		r.responsesQueue <- MessageWrapper{&r.p2, message.Message{Title: message.YourTurn, Payload: nil}}
		r.active = &r.p2
		cellsSlice := r.field.GetAvailableCells(r.getPlayerIdx(r.active))

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

		//Start Timer Here
		r.waitForSyncMsg = message.GoTo

	}
	return true
}

//Выбросить "игрока" из комнаты, поместить в другую (пока не надо трогать)
func (r *Room) ChangeOpponentHandler(m MessageWrapper) bool {

	var secondPlayer *player.Player
	//var thisPlayer *player.Player

	if &r.p1 == m.player {
		//thisPlayer = &r.p1
		secondPlayer = &r.p2
	}
	if &r.p2 == m.player {
		//thisPlayer = &r.p2
		secondPlayer = &r.p1
	}
	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.OpponentLeaves, Payload: nil}}

	return true

}

//Выбросить "пользователя" в главное меню,  connection "игрока" уничтожить
func (r *Room) QuitHandler(m MessageWrapper) bool {

	var secondPlayer *player.Player
	//var thisPlayer *player.Player

	if &r.p1 == m.player {
		secondPlayer = &r.p2
	}
	if &r.p2 == m.player {
		secondPlayer = &r.p1
	}
	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.OpponentLeaves, Payload: nil}}

	return true
}

func (r *Room) getPlayerIdx(p *player.Player) int {
	if &r.p1 == p {
		return 1
	}
	if &r.p2 == p {
		return 2
	}
	logger.Errorf("Player with address %d was not found, couldn't get idx", (*p).UID())
	return 1
}

func (r *Room) CurrentStateHandler(m MessageWrapper) {
	r.responsesQueue <- MessageWrapper{m.player, message.Message{Title: message.CurrentState, Payload: message.GameState{r.field.GetCurrentState()}}}
}

func (r *Room) PackSelectorHandler(m MessageWrapper) bool {
	logger.Info("entered PackSelectorHandler")
	if r.timerToChoosePack.Stop() {
		logger.Info("PackSelectorHandler, Timer is disabled manually")
	} else {
		logger.Info("PackSelectorHandler, Timer is disabled by timeout")
	}

	st, ok := m.msg.Payload.(map[string]interface{})
	if !ok {
		logger.Error("PackSelectorHandler, couldn't cast payload with pack_id to map[string]interface{}")
	}
	packId, ok := st["pack_id"].(float64)
	if !ok {
		logger.Error(`PackSelectorHandler, couldn't find value in map st with key "pack_id" `)
	}

	var secondPlayer *player.Player
	var thisPlayer *player.Player

	if &r.p1 == m.player {
		thisPlayer = &r.p1
		secondPlayer = &r.p2
	}
	if &r.p2 == m.player {
		thisPlayer = &r.p2
		secondPlayer = &r.p1
	}
	//Check if player hasn't answered in time
	if packId == -1 {
		r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.SelectedPack, Payload: message.PackID{PackId: -1}}}
		r.responsesQueue <- MessageWrapper{thisPlayer, message.Message{Title: message.SelectedPack, Payload: message.PackID{PackId: -1}}}

		r.responsesQueue <- MessageWrapper{thisPlayer, message.Message{Title: message.OpponentTurn, Payload: nil}}

		r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.YourTurn, Payload: nil}}
		r.changeTurn()
		r.timerToChoosePack = time.AfterFunc(timeToMove*time.Second, r.ChoosePackTimerFunc)
		r.waitForSyncMsg = message.NotDesiredPack
		return true
	}

	packs := r.field.GetPacksSlice()
	for i, pack := range *packs {
		if pack.ID == uint64(packId) {
			(*packs)[i] = (*packs)[len(*packs)-1] // Replace it with the last one. CAREFUL only works if you have enough elements.
			*packs = (*packs)[:len(*packs)-1]     // Chop off the last one.
			break
		}
		if i+1 == len(*packs) {
			logger.Error("pack with id", packId, "wasn't found in packs slice")
		}
	}
	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.SelectedPack, Payload: message.PackID{PackId: -1}}}
	r.responsesQueue <- MessageWrapper{thisPlayer, message.Message{Title: message.SelectedPack, Payload: message.PackID{PackId: -1}}}

	if len(*packs) == packTotal-2*packsPerPlayer {
		r.waitForSyncMsg = "READY"
		go r.prepareMatch()

		return true
	}

	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.YourTurn, Payload: nil}}
	r.responsesQueue <- MessageWrapper{thisPlayer, message.Message{Title: message.OpponentTurn, Payload: nil}}
	r.changeTurn()
	r.timerToChoosePack = time.AfterFunc(timeToMove*time.Second, r.ChoosePackTimerFunc)
	r.waitForSyncMsg = message.NotDesiredPack
	return true
}

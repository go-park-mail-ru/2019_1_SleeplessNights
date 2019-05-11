package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/event"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player"
)

func (r *Room) MessageHandlerMux(m MessageWrapper) {
	switch m.msg.Title {

	case message.Ready:
		{
			r.ReadyHandler(m)
		}
	case message.GoTo:
		{
			r.GoToHandler(m)
		}
	case message.ClientAnswer:
		{
			r.ClientAnswerHandler(m)
		}
	case message.Leave:
		{
			r.LeaveHandler(m)
		}

	case message.Continue:
		{
			r.ContinueHandler(m)
		}
	case message.ChangeOpponent:
		{
			r.ChangeOpponentHandler(m)
		}
	case message.Quit:
		{
			r.QuitHandler(m)
		}
	case message.State:
		{
			r.CurrentStateHandler(m)
		}
	case message.ThemesRequest:
		{
			r.ThemesRequestHandler(m)
		}
	case message.QuestionsThemesRequest:
		{
			r.QuestionsThemesHandler(m)
		}
	}
}

func (r *Room) ReadyHandler(m MessageWrapper) bool {
	r.mu.Lock()

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

		//Send Available cells to active player (Do it every time, after giving player a turn rights
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.AvailableCells, Payload: cellsSlice}}

	}

	r.mu.Unlock()
	return true
}

func (r *Room) GoToHandler(m MessageWrapper) bool {
	logger.Infof("player UID %d requested GoTo", (*m.player).UID())

	r.mu.Lock()
	var eventSlice []event.Event
	var err error
	var secondPlayer *player.Player

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
		r.mu.Unlock()
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
		}
		if e.Etype == event.WinPrize {
			//Write to DB results of the match
			logger.Info("player", (*r.active).ID(), "Has Won the prize")
			r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.Win, Payload: nil}}
			r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.Loss, Payload: nil}}
			r.waitForSyncMsg = "Leave"
			r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.WannaPlayAgain, Payload: nil}}
			r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.WannaPlayAgain, Payload: nil}}
		}
	}

	r.mu.Unlock()

	return true
}

func (r *Room) ClientAnswerHandler(m MessageWrapper) bool {
	logger.Infof("player UID %d answered to ClientAnswerHandler", (*m.player).UID())
	r.mu.Lock()
	st, ok := m.msg.Payload.(map[string]interface{})
	if !ok {
		logger.Error("ClientAnswerHandler, couldn't cast payload with answer_id to map[string]interface{}")
	}
	answerId, ok := st["answer_id"].(float64)
	if !ok {
		logger.Error(`ClientAnswerHandler, couldn't find value in map st with key "answer_id" `)
	}
	q := r.field.GetRegisterQuestion()

	if !r.field.CheckAnswer(int(answerId)) {
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.YourAnswer, Payload: message.AnswerResult{int(answerId), q.Correct}}}
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

	//Send Available cells to active player (Do it every time, after giving player a turn rights
	r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.AvailableCells, Payload: cellsSlice}}

	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.AvailableCells, Payload: cellsSlice}}

	r.mu.Unlock()
	return true
}

func (r *Room) LeaveHandler(m MessageWrapper) bool {
	r.mu.Lock()

	r.mu.Unlock()
	return true
}

//Оставить комнату с теми же игроками, создать для них новое игровое поле
//Если один из них голосует выйти, то написать об этом другому
func (r *Room) ContinueHandler(m MessageWrapper) bool {
	r.mu.Lock()
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

		//Send Available cells to active player (Do it every time, after giving player a turn rights
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.AvailableCells, Payload: cellsSlice}}

		r.waitForSyncMsg = message.GoTo

	}
	r.mu.Unlock()
	return true
}

//Выбросить "игрока" из комнаты, поместить в другую (пока не надо трогать)
func (r *Room) ChangeOpponentHandler(m MessageWrapper) bool {
	r.mu.Lock()

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

	r.mu.Unlock()
	return true

}

//Выбросить "пользователя" в главное меню,  connection "игрока" уничтожить
func (r *Room) QuitHandler(m MessageWrapper) bool {
	r.mu.Lock()

	var secondPlayer *player.Player
	//var thisPlayer *player.Player

	if &r.p1 == m.player {
		///	thisPlayer=&r.p1
		secondPlayer = &r.p2
	}
	if &r.p2 == m.player {
		///	thisPlayer = &r.p2
		secondPlayer = &r.p1
	}
	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.OpponentLeaves, Payload: nil}}

	r.mu.Unlock()
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

func (r *Room) ThemesRequestHandler(m MessageWrapper) {
	r.responsesQueue <- MessageWrapper{m.player, message.Message{Title: message.Themes, Payload: r.field.GetThemesSlice()}}

}

//Maybe send response to websocket connection without request
func (r *Room) QuestionsThemesHandler(m MessageWrapper) {
	packArray := r.field.GetQuestionsThemes()
	r.responsesQueue <- MessageWrapper{m.player, message.Message{Title: message.QuestionsThemes, Payload: packArray}}
}

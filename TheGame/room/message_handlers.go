package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/event"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/player"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
)

func (r *Room) MessageHandlerMux(m MessageWrapper) {
	switch m.msg.Title {

	case messge.Ready:
		{
			r.ReadyHandler(m)
		}
	case messge.GoTo:
		{
			r.GoToHandler(m)
		}
	case messge.ClientAnswer:
		{
			r.ClientAnswerHandler(m)
		}
	case messge.Leave:
		{
			r.LeaveHandler(m)
		}

	}
}

func (r *Room) ReadyHandler(m MessageWrapper) bool {
	r.mu.Lock()

	if m.player == &r.p1 {
		logger.Info.Println("Игрок 1 Готов")
		r.p1Status = StatusReady
	}
	if m.player == &r.p2 {
		logger.Info.Println("Игрок 2 готов")
		r.p2Status = StatusReady
	}

	if r.p1Status == StatusReady && r.p2Status == StatusReady {
		//After getting ready messages from both players set p1 as active and send messages
		r.waitForSyncMsg = messge.GoTo
		r.active = &r.p1

		logger.Info.Println("Ход игрока 1, ожидание команды GoTo")

		ResponsesQueue <- MessageWrapper{&r.p1, messge.Message{Title: messge.YourTurn, Payload: nil}}
		ResponsesQueue <- MessageWrapper{&r.p2, messge.Message{Title: messge.EnemyTurn, Payload: nil}}
		// Результат работы достаем из канала Events()отсылаем в канал ResponsesQueue
		cellsSlice := r.field.GetAvailableCells(r.getPlayerIdx(r.active))

		//Send Available cells to active player (Do it every time, after giving player a turn rights
		ResponsesQueue <- MessageWrapper{r.active, messge.Message{Title: messge.AvailableCells, Payload: cellsSlice}}

	}

	r.mu.Unlock()
	return true
}

func (r *Room) GoToHandler(m MessageWrapper) bool {
	logger.Info.Printf("player %d requested GoTo", (*m.player).UID())

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
		logger.Error.Println("GoToHandler, called TryMovePLayer, got error", err)
		r.mu.Unlock()
		return false
	}

	for _, e := range eventSlice {
		if e.Etype == event.Info {
			data := e.Edata.(*event.Question)
			ResponsesQueue <- MessageWrapper{r.active, messge.Message{Title: messge.YourQuestion, Payload: data}}
			ResponsesQueue <- MessageWrapper{secondPlayer, messge.Message{Title: messge.EnemyQuestion, Payload: data}}

		}
	}

	r.mu.Unlock()

	return true
}

func (r *Room) ClientAnswerHandler(m MessageWrapper) bool {
	logger.Info.Printf("player %d answered to ClientAnswerHandler", (*m.player).UID())
	r.mu.Lock()
	answerId := m.msg.Payload.(*messge.Answer).AnswerId
	if !r.field.CheckAnswer(answerId) {
		ResponsesQueue <- MessageWrapper{r.active, messge.Message{Title: messge.Incorrect, Payload: nil}}
	}
	ResponsesQueue <- MessageWrapper{r.active, messge.Message{Title: messge.Correct, Payload: nil}}

	//Смена хода после ответа игрока
	ResponsesQueue <- MessageWrapper{r.active, messge.Message{Title: messge.EnemyTurn, Payload: nil}}
	r.changeTurn()
	r.waitForSyncMsg = "GoTo"
	ResponsesQueue <- MessageWrapper{r.active, messge.Message{Title: messge.YourTurn, Payload: nil}}

	r.mu.Unlock()
	return true
}

func (r *Room) LeaveHandler(m MessageWrapper) bool {

	r.mu.Lock()

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
	logger.Error.Println("Player with address %d was not found, couldn't get idx", p)
	return 1
}

package room

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/event"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"time"
)

type Pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

const (
	timeToWait = 2
)

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
		r.active = &r.p1
		go r.startGameProcess()

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
				Time       int
			}{
				CellsSlice: cells,
				Time:       timeToMove,
			}
			//Send Available cells to active player (Do it every time, after giving player a turn rights
			r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.AvailableCells, Payload: payload}}
			r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.AvailableCells, Payload: payload}}
			r.timerToMove = time.AfterFunc(time.Duration(timeToAnswer)*time.Second, r.GoToTimerFunc)
			r.waitForSyncMsg = message.GoTo
			return true
		} else {
			logger.Error("Unexpected condition")
		}
	}
	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.SelectedCell, Payload: message.Coordinates{nextX, nextY}}}
	r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.SelectedCell, Payload: message.Coordinates{nextX, nextY}}}

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
			r.timerToAnswer = time.AfterFunc(time.Duration(timeToAnswer)*time.Second, r.AnswerTimerFunc)
		}
		if e.Etype == event.WinPrize {
			//Write to DB results of the
			logger.Info("player", (*r.active).ID(), "Has Won the prize")
			r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.Win, Payload: nil}}
			r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.Loss, Payload: nil}}
			if secondPlayer == nil {
				logger.Error("secondPlayer, Attempt nil dereference ")
			}

			var winnerRating uint64
			var loserRating uint64
			idx := r.getPlayerIdx(r.active)
			if idx == 1 {
				winnerRating = r.p1Rating
				loserRating = r.p2Rating
			} else {
				winnerRating = r.p2Rating
				loserRating = r.p1Rating
			}
			if secondPlayer != nil {
				results := services.MatchResults{
					Winner:       (*r.active).UID(),
					Loser:        (*secondPlayer).UID(),
					WinnerRating: winnerRating,
					LoserRating:  loserRating,
				}
				_, err := userManager.UpdateStats(context.Background(), &results)
				if err != nil {
					logger.Error("Failed to update Match Statistics:", err)
				}
				(*secondPlayer).Close()
				logger.Info("WinPrize, second Close() called")
			}
			r.KillMePleaseFlag = true
			if r.timerToAnswer != nil {
				r.timerToAnswer.Stop()
			}
			if r.timerToMove != nil {
				r.timerToMove.Stop()
			}
			if r.timerToChoosePack != nil {
				r.timerToChoosePack.Stop()
			}

			(*r.active).Close()
			logger.Info("WinPrize, active Close() called")
			logger.Info("Won The Prize, room is to be deleted, r.KillMePleaseFlag = true")

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
		} else {
			idx := r.getPlayerIdx(r.active)
			if idx == 1 {
				r.p1Rating += 1
			} else {
				r.p2Rating += 1
			}
		}
	} else {
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.YourAnswer, Payload: message.AnswerResult{int(answerId), q.Correct}}}
		r.field.Move(r.getPlayerIdx(r.active))
	}
	var secondPlayer *player.Player

	if &r.p1 == r.active {
		secondPlayer = &r.p2
	} else {
		secondPlayer = &r.p1
	}

	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.OpponentAnswer, Payload: message.AnswerResult{int(answerId), q.Correct}}}

	if playerHasNoMoves {
		r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.Loss, Payload: nil}}
		r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.Win, Payload: nil}}

		var winnerRating uint64
		var loserRating uint64
		idx := r.getPlayerIdx(r.active)
		if idx == 1 {
			winnerRating = r.p1Rating
			loserRating = r.p2Rating
		} else {
			winnerRating = r.p2Rating
			loserRating = r.p1Rating
		}

		if secondPlayer == nil {
			logger.Error("Attempt nill dereference of secondPlayer pointer ")
			return true
		}

		results := services.MatchResults{
			Winner:       (*secondPlayer).UID(),
			Loser:        (*r.active).UID(),
			WinnerRating: winnerRating,
			LoserRating:  loserRating,
		}
		_, err := userManager.UpdateStats(context.Background(), &results)
		if err != nil {
			logger.Error("Failed to update Match Statistics:", err)
		}
		r.KillMePleaseFlag = true

		if r.timerToAnswer != nil {
			r.timerToAnswer.Stop()
		}
		if r.timerToMove != nil {
			r.timerToMove.Stop()
		}
		if r.timerToChoosePack != nil {
			r.timerToChoosePack.Stop()
		}

		logger.Info("No Moves Left, r.KillMePleaseFlag = true")
		r.p1.Close()
		r.p2.Close()
		logger.Info("No Moves Left, Close channels")

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
				Time       int
			}{
				CellsSlice: cells,
				Time:       timeToMove,
			}
			//Send Available cells to active player (Do it every time, after giving player a turn rights
			r.responsesQueue <- MessageWrapper{r.active, message.Message{Title: message.AvailableCells, Payload: payload}}
			r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.AvailableCells, Payload: payload}}
			r.timerToMove = time.AfterFunc(time.Duration(timeToMove)*time.Second, r.GoToTimerFunc)

		} else {
			logger.Error("Unexpected condition")
		}
	}
	return true
}

func (r *Room) LeaveHandler(m MessageWrapper) bool {

	var leaverPlayer *player.Player
	leaverPlayer = m.player
	var stayerPlayer *player.Player
	leaver_idx := ""

	if r.timerToAnswer != nil {
		r.timerToAnswer.Stop()
	}
	if r.timerToMove != nil {
		r.timerToMove.Stop()
	}
	if r.timerToChoosePack != nil {
		r.timerToChoosePack.Stop()
	}
	close((*leaverPlayer).Subscribe())
	if &r.p1 == leaverPlayer {
		leaver_idx = "1"
		stayerPlayer = &r.p2

		if stayerPlayer != nil {
			r.responsesQueue <- MessageWrapper{stayerPlayer, message.Message{message.Leave, "Player2 left the game"}}
		}
	} else {
		leaver_idx = "2"
		stayerPlayer = &r.p1
		if stayerPlayer != nil {
			r.responsesQueue <- MessageWrapper{stayerPlayer, message.Message{message.Leave, "Player1 left the game"}}
		}
	}
	if stayerPlayer != nil {
		r.responsesQueue <- MessageWrapper{stayerPlayer, message.Message{Title: message.Win, Payload: nil}}

		logger.Info("Leave Handler, Player" + leaver_idx + " ID " + fmt.Sprint((*leaverPlayer).ID()) + "Closed Connection")
		var winnerRating uint64
		var loserRating uint64

		idx := r.getPlayerIdx(stayerPlayer)
		if idx == 1 {
			winnerRating = r.p1Rating
			loserRating = r.p2Rating
		} else {
			winnerRating = r.p2Rating
			loserRating = r.p1Rating
		}

		results := services.MatchResults{
			Winner:       (*stayerPlayer).UID(),
			Loser:        (*leaverPlayer).UID(),
			WinnerRating: winnerRating,
			LoserRating:  loserRating,
		}
		_, err := userManager.UpdateStats(context.Background(), &results)
		if err != nil {
			logger.Error("Failed to update Match Statistics:", err)
		}
		r.KillMePleaseFlag = true

		if &r.p1 == leaverPlayer {
			r.p2.Close()
		}

		if &r.p2 == leaverPlayer {
			r.p1.Close()
		}

		logger.Info("Player Leaves, r.KillMePleaseFlag = true")

	} else {
		logger.Info("Player left empty room, room is to be deleted")
	}
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

		r.responsesQueue <- MessageWrapper{thisPlayer, message.Message{Title: message.OpponentTurn, Payload: nil}}

		r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.YourTurn, Payload: nil}}
		r.changeTurn()
		r.timerToChoosePack = time.AfterFunc(time.Duration(timeToMove)*time.Second, r.ChoosePackTimerFunc)
		r.waitForSyncMsg = message.NotDesiredPack
		return true
	}
	logger.Info("player chosen pack_ID", packId)
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
	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.SelectedPack, Payload: message.PackID{PackId: int64(packId)}}}

	if len(*packs) == packTotal-2*packsPerPlayer {
		r.active = &r.p1
		go r.prepareMatch()
		return true
	}

	r.responsesQueue <- MessageWrapper{secondPlayer, message.Message{Title: message.YourTurn, Payload: nil}}
	r.responsesQueue <- MessageWrapper{thisPlayer, message.Message{Title: message.OpponentTurn, Payload: nil}}
	r.changeTurn()
	r.timerToChoosePack = time.AfterFunc(time.Duration(timeToMove)*time.Second, r.ChoosePackTimerFunc)
	r.waitForSyncMsg = message.NotDesiredPack
	return true
}

package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
)

//Проверка Уместности сообщения ( на уровне комнаты)
func (r *Room) isSyncValid(wm MessageWrapper) (isValid bool) {
	r.mu.Lock()
	if wm.msg.Title == message.Leave {
		isValid = true
		r.mu.Unlock()
		return
	}
	if wm.msg.Title == message.ChangeOpponent || wm.msg.Title == message.Quit || wm.msg.Title == message.Continue {
		isValid = true
		r.mu.Unlock()
		return
	}
	if wm.msg.Title == message.State {
		isValid = true
		r.mu.Unlock()
		return
	}

	if wm.msg.Title == message.NotDesiredPack && r.active == wm.player {
		isValid = true
		r.mu.Unlock()
		return
	}

	if wm.msg.Title == message.ThemesRequest {
		isValid = true
		r.mu.Unlock()
		return
	}

	if wm.player != r.active && (wm.msg.Title != message.Ready) {
		logger.Error("isSync Player addr error")
		isValid = false
		r.mu.Unlock()
		return
	}
	if r.waitForSyncMsg != wm.msg.Title {
		logger.Error("isSync title error")
		isValid = false
		r.mu.Unlock()
		return
	}
	isValid = true
	r.mu.Unlock()
	return
}

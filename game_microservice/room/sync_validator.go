package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
)

//Проверка Уместности сообщения ( на уровне комнаты)
func (r *Room) isSyncValid(wm MessageWrapper) (isValid bool) {

	if wm.msg.Title == message.Leave {
		isValid = true
		return
	}
	if wm.msg.Title == message.State {
		isValid = true

		return
	}

	if wm.msg.Title == message.NotDesiredPack {
		isValid = true
		return
	}
	if wm.msg.Title == message.ThemesRequest {
		isValid = true
		return
	}
	if wm.player != r.active && (wm.msg.Title != message.Ready) {
		logger.Error("isSync Player addr error")

		isValid = false
		return
	}
	if r.waitForSyncMsg != wm.msg.Title {
		logger.Error("isSync title error")
		isValid = false
		return
	}
	isValid = true
	return
}

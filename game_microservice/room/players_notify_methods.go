package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
)

func (r *Room) notifyP1(msg message.Message) (err error) {
	err = r.p1.Send(msg)
	if err != nil {
		logger.Error("Failed to send Message to P1", err)
	}
	return
}

func (r *Room) notifyP2(msg message.Message) (err error) {
	err = r.p2.Send(msg)
	if err != nil {
		logger.Error("Failed to send Message to P2", err)
	}
	return
}

func (r *Room) notifyAll(msg message.Message) (err error) {
	err = r.notifyP1(msg)
	if err != nil {
		return
	}
	err = r.notifyP2(msg)
	if err != nil {
		return
	}
	return nil
}

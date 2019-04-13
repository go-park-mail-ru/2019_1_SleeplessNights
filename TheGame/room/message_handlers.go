package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
)

func (r *Room) MessageHandlerMux(m MessageWrapper) {

	switch m.msg.Title {

	case messge.Ready:
		{
			//"Ready" Message is awaited in r.startmatch, maybe should add room.readyflags for each player
		}

	case messge.GoTo:
		{

		}
	case messge.ClientAnswer:
		{

		}
	case messge.Leave:
		{

		}

	}
}

func (r *Room) ReadyHandler(msg messge.Message) bool {
	return true
}

func (r *Room) GoToHandler(msg messge.Message) bool {
	return true
}

func (r *Room) ClientAnswerHandler(msg messge.Message) bool {
	return true
}

func (r *Room) LeaveHandler(message messge.Message) bool {
	return true
}

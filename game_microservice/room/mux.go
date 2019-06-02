package room

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"

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
	case message.State:
		{
			r.CurrentStateHandler(m)
		}

	case message.NotDesiredPack:
		{
			r.PackSelectorHandler(m)
		}
	default:
		{
			logger.Error("MessageHandlerMux, unknown title")
		}
	}
}

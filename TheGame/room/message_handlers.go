package room

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"

func (r *Room) MessageHandlerMux(m messge.Message) {
	switch m.Title {
	case "COMMAND":
		{
			if m.CommandName == messge.CommandMove {
				//player requests moving to x,y cell
				//If cell is available, then send him a question
				if r.MoveHandler(m) {

				}
			}
			if m.CommandName == messge.CommandAnswer {
				//Check if answer is correct
				// Do something
				if AnswerHandler(m) {
				}
			}
			//...

		}

	case "INFO":
		{

		}

		//....
	}

}

func (r *Room) MoveHandler(message messge.Message) bool {
	payload, ok := message.Payload.(*messge.MoveRequest)
	if !ok {
		return false
	}

	return true
}

func AnswerHandler(message messge.Message) bool {
	payload, ok := message.Payload.(*messge.Answer)
	if !ok {
		return false
	}
	return true
}

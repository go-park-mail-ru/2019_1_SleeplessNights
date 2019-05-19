package room

import (
	"time"
)

func (r *Room) StartRequestsListening() {

	logger.Info("started listening to messages channels")
	r.requestsQueue = make(chan MessageWrapper, channelCapacity)
	r.responsesQueue = make(chan MessageWrapper, channelCapacity)

	go func() {
		for msgP1 := range r.p1.Subscribe() {
			logger.Info("got message from P1", msgP1)
			r.requestsQueue <- MessageWrapper{&r.p1, msgP1}
		}
	}()

	go func() {
		for msgP2 := range r.p2.Subscribe() {
			logger.Info("got message from P2", msgP2)
			r.requestsQueue <- MessageWrapper{&r.p2, msgP2}
		}
	}()

	go func() {
		for msg := range r.requestsQueue {
			//Проверка структуры сообщения
			logger.Info("Got Message from client")
			if !msg.msg.IsValid() {
				logger.Error("Got message with invalid structure")
				continue
			}

			if !r.isSyncValid(msg) {
				logger.Warningf("Got SyncInvalid message of type %s from player UID %d, expected %s from player %d",
					msg.msg.Title, (*msg.player).UID(), r.waitForSyncMsg, r.active)
				continue
			}
			logger.Info("Message entered mux")
			r.MessageHandlerMux(msg)
		}
	}()

}

func (r *Room) StartResponsesSender() {

	logger.Info("started listening to response channel")
	go func() {
		for serverResponse := range r.responsesQueue {
			logger.Info("Got message to Send recepient: UID", (*serverResponse.player).UID(), "Message:", serverResponse.msg)
			err := (*serverResponse.player).Send(serverResponse.msg)
			if err != nil {
				logger.Error("responseQueue: error trying to send response to player", err)
			}
		}
		time.Sleep(responseInterval * time.Millisecond)
	}()
}

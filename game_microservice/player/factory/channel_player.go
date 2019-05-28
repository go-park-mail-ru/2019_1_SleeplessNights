package factory

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
)

type ChannelPlayerLogic func(id uint64, in, out *chan message.Message, args ...interface{})()//id - просто идентификатор для логов

type channelPlayer struct {
	//С помощью этой структуры будем делать игрока по вебсокету на фабрике
	//Она уже реализует интерфейс Player, поэтому нам нужно будет просто сделатьинстанс этой структуры
	work ChannelPlayerLogic
	id   uint64
	in   chan message.Message
	out  chan message.Message
}

func (chanPlayer *channelPlayer) Send(msg message.Message) (err error) {
	chanPlayer.out <- msg
	return nil
}

func (chanPlayer *channelPlayer) Subscribe() chan message.Message {
	return chanPlayer.in
}

func (chanPlayer *channelPlayer) ID() uint64 {
	return chanPlayer.id
}

func (chanPlayer *channelPlayer) UID() uint64 {
	return 0
}

func (chanPlayer *channelPlayer) Close() {
	logger.Infof("Player %d closed the connection", chanPlayer.id)
	chanPlayer.in <- message.Message{Title: message.Leave}
	close(chanPlayer.in)
	close(chanPlayer.out)
}

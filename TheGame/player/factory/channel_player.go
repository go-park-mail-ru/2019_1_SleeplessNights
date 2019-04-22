package factory

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
)

type channelPlayerLogic func(args ...interface{})()

type channelPlayer struct {
	//С помощью этой структуры будем делать игрока по вебсокету на фабрике
	//Она уже реализует интерфейс Player, поэтому нам нужно будет просто сделатьинстанс этой структуры
	work channelPlayerLogic
	id   uint64
	in   chan messge.Message
	out  chan messge.Message
}

func (chanPlayer *channelPlayer) Send(msg messge.Message) (err error) {
	chanPlayer.out <- msg
	return nil
}

func (chanPlayer *channelPlayer) Subscribe() chan messge.Message {
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
	chanPlayer.in <- messge.Message{Title: messge.Leave}
	close(chanPlayer.in)
	close(chanPlayer.out)
}

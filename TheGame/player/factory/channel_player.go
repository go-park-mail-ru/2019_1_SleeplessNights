package factory

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
)

type channelPlayer struct {
	//TODO develop Close() method
	//С помощью этой структуры будем делать игрока по вебсокету на фабрике
	//Она уже реализует интерфейс Player, поэтому нам нужно будет просто сделатьинстанс этой структуры
	id    uint64
	uid   uint64
	In    chan messge.Message
	Outer chan messge.Message
}

func (chanPlayer *channelPlayer) StartListen() {
	logger.Info("ChannelPlayer Started Listening")

	for msg := range chanPlayer.Outer {
		logger.Info("Got from Outer channel", msg)
		logger.Info("Passed message from outer to Inner Channel")
		chanPlayer.In <- msg
	}
}

func (chanPlayer *channelPlayer) Send(msg messge.Message) (err error) {
	chanPlayer.Outer <- msg
	return nil
}

func (chanPlayer *channelPlayer) Subscribe() chan messge.Message {
	return chanPlayer.In
}

func (chanPlayer *channelPlayer) ID() uint64 {
	return chanPlayer.id
}

func (chanPlayer *channelPlayer) UID() uint64 {
	return chanPlayer.uid
}

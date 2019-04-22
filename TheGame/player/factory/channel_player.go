package factory

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
)

type channelPlayer struct {
	//С помощью этой структуры будем делать игрока по вебсокету на фабрике
	//Она уже реализует интерфейс Player, поэтому нам нужно будет просто сделатьинстанс этой структуры
	id  uint64
	In  chan messge.Message
	Out chan messge.Message
}

func (chanPlayer *channelPlayer) Send(msg messge.Message) (err error) {
	chanPlayer.Out <- msg
	return nil
}

func (chanPlayer *channelPlayer) Subscribe() chan messge.Message {
	return chanPlayer.In
}

func (chanPlayer *channelPlayer) ID() uint64 {
	return chanPlayer.id
}

func (chanPlayer *channelPlayer) UID() uint64 {
	return 0
}

func (chanPlayer *channelPlayer) Close() {
	logger.Infof("Player %d closed the connection", chanPlayer.id)
	chanPlayer.In <- messge.Message{Title: messge.Leave}
	close(chanPlayer.In)
	close(chanPlayer.Out)
}

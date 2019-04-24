package player

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/messge"
)

type Player interface {
	//Игрок в нашей реализации - любое однозначно идентифицируемое дуплексное соединение
	//С помощью метода Send мы будем слать сообщения игроку
	//С помощью метода Subscribe, мы будем получать канал, в который будут приходить сообщения от игрока
	//Метод ID даёт нам некоторый суррогатный ключ, по которому мы можем однозначно идентифицировать игрока
	//ВАЖНО! о природе ключа мы ничего не знаем
	Send(msg messge.Message) error
	Subscribe() chan messge.Message
	ID() uint64
	UID() uint64
}

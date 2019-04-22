package factory

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/player"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/gorilla/websocket"
	"sync/atomic"
)

//В этом файле находится всё, что связано с внутренней реализацией фабрики игроков
//Этот объект должен принимать на вход какое-то соединение и выдавать на выходе (например websocket conn)
//Объект, реализующий интерфейс Player, то есть понятный для всего остального микросервиса

//В ReadMe я написал, что наша фабрика - синглтон
//Логика такого решения следующая - зачем нам игроки, которые были сделаны в разных местах? Пусть у них будет единый источник
//Всё это значит, что на весь проект у нас есть ТОЛЬКО ОДНА фабрика игроков

var logger *log.Logger

func init() {
	logger = log.GetLogger("PlayerFactory")
}

var factory *playerFactory //Вот она
//Обратите внимание, что она неэксортируемая, то есть мы не можем из другого пакета напрямую взять это значение
//и сделать какую-нибудь грязь, в том числе и скопировать

type playerFactory struct {
	//Здесь находится описание структуры нашей фабрики
	idSource uint64 //Счётчик, по которому будем выдавать игрокам ID
}

func init() {
	//При компиляции инициализируем нашу фабрику
	factory = &playerFactory{
		idSource: 0,
	}
}

func GetInstance() *playerFactory {
	//Функция экспортируемая, т.е. теперь в других пакетах мы можем получать указатель на нашу фабрику
	//и работать во всей программе со множеством указателей на один конкретный инстанс фабрики
	return factory
}

func (pf *playerFactory) BuildWebsocketPlayer(conn *websocket.Conn, uid uint64) player.Player {
	//Метод построения игрока из вебсокет соединения
	wsPlayer := websocketPlayer{
		id:   atomic.AddUint64(&pf.idSource, 1), //Атомик необходим для обеспечения потокобезопасности
		conn: conn,
		uid:  uid,
		in:   make(chan messge.Message, 1),
	}
	go wsPlayer.StartListen()
	logger.Info("wsPlayer started listening", uid)
	return &wsPlayer
}

func (pf *playerFactory) BuildChannelPlayer(jobToDo channelPlayerLogic, args ...interface{}) player.Player {
	//Метод построения игрока из вебсокет соединения
	chanPlayer := channelPlayer{
		work: jobToDo,
		id:   atomic.AddUint64(&pf.idSource, 1), //Атомик необходим для обеспечения потокобезопасности
		in:   make(chan messge.Message, 1),
		out:  make(chan messge.Message, 1),
	}
	go chanPlayer.work(args...)
	logger.Info("chanPlayer started working", chanPlayer.id)
	return &chanPlayer
}

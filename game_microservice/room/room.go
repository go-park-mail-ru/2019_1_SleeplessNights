package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game_field"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player"
	"sync"
	"time"
)

//Комната - тот объект, который инкапсулирует в себе всю рботу с игровыми механиками
//Здесь нам нужно внутри игрового цикла слать получать сообщения игроков о том, какие ходы они делают,
//прощитывать новые игровые ситуации и отправлять игрокам сообщения об изменении игровой ситуации

//Из неигровых задач комната должна уметь
//* Собираться и пересобираться, не выкидывая игроков, если оба решили сыграть ещё партию вместе или это лобби
//* Поддерживать обработку отвалившегося игрока

const (
	responseInterval = 500
	channelCapacity  = 50
	packTotal        = 10
	packsPerPlayer   = 2
	timeToAnswer     = 20
	timeToMove       = 20
)

const (
	StatusJoined = iota
	StatusReady
	StatusLeft
	StatusSelectedPacks
	StatusWannaContinue
	StatusWannaChangeOpponent
	StatusWannaQuit
)

type MessageWrapper struct {
	player *player.Player
	msg    message.Message
}

type Room struct {
	//TODO develop Close() method
	//Channel to exchange event messages between Room and GameField
	requestsQueue  chan MessageWrapper
	responsesQueue chan MessageWrapper
	p1             player.Player
	p2             player.Player
	p1Status       int
	p2Status       int
	active         *player.Player
	mu             sync.Mutex //Добавление игрока в комнату - конкурентная операция, поэтому нужен мьютекс
	field          game_field.GameField
	waitForSyncMsg string
	timerToAnswer  *time.Timer
	timerToMove    *time.Timer
	syncChan       chan bool
	//Если не знаете, что это такое, то погуглите (для любого языка), об этом написано много, но, обычно, довольно сложно
	//Если по-простому, то это типа стоп-сигнала для всех остальных потоков, который можно включить,
	//сделать всё, что нужно, пока тебе никто не мешает, и выключить обратно
}

func (r *Room) TryJoin(p player.Player) (success bool) {
	//Здесь нам нужно под мьютексом проверить наличие свободных мест. Варианты:
	//1. 2 места свободно -> занимаем первое место
	//2. Свободно 1 место -> занимаем место, поднимаем флаг недоступности комнаты, начинаем игровой процесс
	logger.Infof("player with UID %d entered Try Join", p.UID())

	found := false

	if r.p1 == nil {
		r.p1 = p
		r.p1Status = StatusJoined
		logger.Infof("Player  with UID %d is now p1 in room", p.UID())
		err := r.notifyP1(message.Message{Title: "CONNECTED", Payload: "you've been added to room"})
		if err != nil {
			logger.Error("Failed to notify player ", p.UID())
		}
		found = true
	} else if r.p2 == nil {
		r.p2 = p
		r.p1Status = StatusJoined
		err := r.notifyP2(message.Message{Title: "CONNECTED", Payload: "you've been added to room"})
		if err != nil {
			logger.Error("Failed to notify player ", p.UID())
		}
		logger.Infof("Player UID %d is now p2 in room", p.UID())

		found = true
	}

	if r.p1 != nil && r.p2 != nil {
		logger.Infof("All players joined the Room, p1 UID: %d, p2 UID: %d", r.p1.UID(), r.p2.UID())
		logger.Info("Started Listening Messages ")
		r.active = &r.p1
		r.StartRequestsListening()
		r.StartResponsesSender()
		go r.startGameProcess()
	}
	return found
}

func (r *Room) grantGodMod(p player.Player, token []byte) {
	//РЕАЛИЗОВЫВАТЬ ПОСЛЕДНЕЙ
	//Это чисто техническая процедура, она нужна не для реальных игроков, а, в основном, для ботов, которые должны знать
	//правильный ответ, чтобы отвечать верно более чем на 25% вопросов
	//Принцип работы следующий:
	//1. У игрока есть токен на получение всех овтветов (бот будет запрашивать его у сервера, из другого пакета)
	//2. Игрок отправляет сообщение с запросом на все ответы и своим токеном
	//3. Получив сообщение, комната запускает эту функцию
	//4. Здесь мы проверяем валидность токена, и возвращаем в сообщении игроку матрицу правильных ответов
	//5. ВАЖНО! Конретно это сообщение надо отправлять напрямую конкретному игроку, а не через notify
	//TODO develop
}

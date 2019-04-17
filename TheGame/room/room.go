package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/game_field"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/player"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
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
)

const (
	StatusJoined = iota
	StatusReady
	StatusLeft
)

type MessageWrapper struct {
	player *player.Player
	msg    messge.Message
}

//При таком обьявлении каналы будут общими для всей программы- это плохо, пихнуть их в структуру комнаты

var RequestsQueue chan MessageWrapper
var ResponsesQueue chan MessageWrapper

type Room struct {
	//Channel to exchange event messages between Room and GameField
	p1             player.Player
	p2             player.Player
	p1Status       int
	p2Status       int
	active         *player.Player
	field          game_field.GameField
	waitForSyncMsg string
	mu             sync.Mutex //Добавление игрока в комнату - конкурентная операция, поэтому нужен мьютекс
	//Если не знаете, что это такое, то погуглите (для любого языка), об этом написано много, но, обычно, довольно сложно
	//Если по-простому, то это типа стоп-сигнала для всех остальных потоков, который можно включить,
	//сделать всё, что нужно, пока тебе никто не мешает, и выключить обратно
}

func init() {

}

func (r *Room) TryJoin(p player.Player) (success bool) {
	//Здесь нам нужно под мьютексом проверить наличие свободных мест. Варианты:
	//1. 2 места свободно -> занимаем первое место
	//2. Свободно 1 место -> занимаем место, поднимаем флаг недоступности комнаты, начинаем игровой процесс
	logger.Info.Printf("player %d entered Try Join", p.UID())
	r.mu.Lock()
	found := false

	if r.p1 == nil {
		r.p1 = p
		r.p1Status = StatusJoined
		logger.Info.Printf("Player %d is now p1 in room", p.UID())
		found = true
	} else if r.p2 == nil {
		r.p2 = p
		r.p1Status = StatusJoined
		logger.Info.Printf("Player %d is now p2 in room", p.UID())

		found = true
	}

	if r.p1 != nil && r.p2 != nil {
		logger.Info.Println("All players joined the game, p1: %d, p2: %d", r.p1.UID(), r.p1.UID())
		//TODO Prepare Match
		//TODO Then run buildEnv after PrepareMatch
		// In build Env составление и доставание даннных для вопросов
		go r.startMatch()

	}

	r.mu.Unlock()
	return found
}

func (r *Room) buildEnv() {

	//TODO Get questions from database

	//TODO CAll GameField.build

	//Процедура должна пересоздавать игровое поле, запрашивать новый список тем из БД и готовить комнату к новой игре
	//При этом она должна уметь работать асинхронно и не выбрасывать пользователей из комнаты во время работы
	//TODO develop

}

func (r *Room) notifyP1(msg messge.Message) (err error) {
	err = r.p1.Send(msg)
	if err != nil {
		logger.Error.Println("Failed to send Message to P1", err)
	}
	return
}

func (r *Room) notifyP2(msg messge.Message) (err error) {
	err = r.p2.Send(msg)
	if err != nil {
		logger.Error.Println("Failed to send Message to P2", err)
	}
	return
}

func (r *Room) notifyAll(msg messge.Message) (err error) {
	err = r.notifyP1(msg)
	if err != nil {
		return
	}

	err = r.notifyP2(msg)
	if err != nil {
		return
	}
	return nil
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

// TODO PREPAREMATCH AND BUILD ENV (simultaneously (optional), then wait them both to work out, use with WaitGroup )
//

func (r *Room) PrepareRoom() {
	logger.Info.Println("Entered Prepare Match")

}

func (r *Room) startMatch() {
	//   Эта процедура запускает игровой процесс
	//   Здесь мы будем слушать все сообщения пользователей асинхронно и складывать их в очередь для обработки
	//   В цикле мы будем обрабатывать все входные сообщения, выполнять нашу бизнес логику (менять значение таймера,
	//   давать пользователям вопросы и т.д.)
	//   Все сообщения мы будем складывать в очередь на отправку и отправлять всю очередь каждые 0.5 сек
	//   (цифра примерная, может поменяться и должна быть вынесена в костанту)
	//TODO develop

	// Call Prepare Room

	logger.Info.Println("StartMatch : Game process has started p1:%s, p2: %s", r.p1.UID(), r.p2.UID())
	r.buildEnv()

	p1Chan := r.p1.Subscribe()
	p2Chan := r.p2.Subscribe()

	err := r.notifyAll(messge.Message{Title: messge.StartGame, Payload: nil})
	if err != nil {
		logger.Error.Printf("Failed to notify all players %s", err)
	}
	logger.Info.Println("Игрокам Отправлены StartGame")

	//Read Messages from Players
	//Moved message receive conditions to Requests handler
	go func() {
		for msgP1 := range p1Chan {
			logger.Info.Println("got message from P1", msgP1)
			RequestsQueue <- MessageWrapper{&r.p1, msgP1}
		}
	}()

	go func() {
		for msgP2 := range p2Chan {
			logger.Info.Println("got message from P2", msgP2)
			RequestsQueue <- MessageWrapper{&r.p2, msgP2}
		}
	}()
	//Channel to Write Server messages to the player1/player2

	go func() {
		for serverResponse := range ResponsesQueue {
			logger.Info.Println("Got message to Send", serverResponse)
			err := (*serverResponse.player).Send(serverResponse.msg)
			if err != nil {
				logger.Error.Println("responseQueue: error trying to send response to player", err)
			}
		}
		time.Sleep(responseInterval * time.Millisecond)
	}()

	go func() {
		for msg := range RequestsQueue {
			//Проверка структуры сообщения
			if !msg.msg.IsValid() {
				logger.Error.Println("Message to send has invalid structure")
				continue
			}

			if r.isSyncValid(msg) {
				logger.Warning.Printf("Got message of type %s from player %d, expected %s from player %d",
					msg.msg.Title, msg.player, r.waitForSyncMsg, r.active)
				continue
			}

			r.MessageHandlerMux(msg)

		}
	}()
}

//Проверка Уместности сообщения ( на уровне комнаты)
func (r *Room) isSyncValid(wm MessageWrapper) (isValid bool) {
	r.mu.Lock()
	if wm.player != r.active {
		isValid = false
	}
	if r.waitForSyncMsg != wm.msg.Title {
		isValid = false
	}
	isValid = true
	r.mu.Unlock()
	return
}

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
)

type Room struct {
	//Channel to exchange event messages between Room and GameField
	p1     player.Player
	p2     player.Player
	active *player.Player
	field  game_field.GameField

	mu sync.Mutex //Добавление игрока в комнату - конкурентная операция, поэтому нужен мьютекс
	//Если не знаете, что это такое, то погуглите (для любого языка), об этом написано много, но, обычно, довольно сложно
	//Если по-простому, то это типа стоп-сигнала для всех остальных потоков, который можно включить,
	//сделать всё, что нужно, пока тебе никто не мешает, и выключить обратно
}

func (r *Room) TryJoin(p player.Player) (success bool) {
	//Здесь нам нужно под мьютексом проверить наличие свободных мест. Варианты:
	//1. 2 места свободно -> занимаем первое место

	//2. Свободно 1 место -> занимаем место, поднимаем флаг недоступности комнаты, начинаем игровой процесс

	r.mu.Lock()
	found := false

	if r.p1 == nil {
		r.p1 = p
		found = true
	} else if r.p2 == nil {
		r.p2 = p
		found = true
	}

	if r.p1 != nil && r.p2 != nil {
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

func (r *Room) notifyAll(msg messge.Message) {
	//Процедуры должна отправить сообщение всем игрокам комнаты
	//TODO develop
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

func (r *Room) startMatch() {

	//   Эта процедура запускает игровой процесс
	// v Здесь мы будем слушать все сообщения пользователей асинхронно и складывать их в очередь для обработки
	//   В цикле мы будем обрабатывать все входные сообщения, выполнять нашу бизнес логику (менять значение таймера,
	//   давать пользователям вопросы и т.д.)
	//   Все сообщения мы будем складывать в очередь на отправку и отправлять всю очередь каждые 0.5 сек
	//   (цифра примерная, может поменяться и должна быть вынесена в костанту)
	//TODO develop

	r.buildEnv()
	//TODO Notify all about start of game

	RequestsQueue := make(chan messge.Message, 50)
	ResponsesQueue := make(chan messge.Message, 50)

	p1Chan := r.p1.Subscribe()
	p2Chan := r.p2.Subscribe()

	//Read Messages from Players
	go func() {
		for msgP1 := range p1Chan {
			logger.Info.Println("got message from P1", msgP1)
			RequestsQueue <- msgP1 //TODO only if msgP1 is valid and (p1 is active or msgP1 is ASYNC)
		}
	}()

	go func() {
		for msgP2 := range p2Chan {
			logger.Info.Println("got message from P2", msgP2)
			RequestsQueue <- msgP2 //TODO only if msgP2 is valid and (p2 is active or msgP2 is ASYNC)
		}
	}()

	go func() {
		for {
			for serverResponse := range ResponsesQueue {
				logger.Info.Println("Got message to Send", serverResponse)

				//TODO Send Message Here
			}
			time.Sleep(responseInterval * time.Millisecond)
		}
	}()

	go func() {
		for event := range r.field.Out {
			logger.Info.Println("Got event from game_field", event)
			//TODO Handle Event Here
		}
	}()

	//TODO Handler players requests
	for msg := range RequestsQueue {
		r.MessageHandlerMux(msg)
	}

}

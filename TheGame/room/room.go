package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/game_field"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/player"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"sync"
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
	StatusWannaContinue
	StatusWannaChangeOpponent
	StatusWannaQuit
)

type MessageWrapper struct {
	player *player.Player
	msg    messge.Message
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
	field          game_field.GameField
	waitForSyncMsg string
	mu             sync.Mutex //Добавление игрока в комнату - конкурентная операция, поэтому нужен мьютекс
	//Если не знаете, что это такое, то погуглите (для любого языка), об этом написано много, но, обычно, довольно сложно
	//Если по-простому, то это типа стоп-сигнала для всех остальных потоков, который можно включить,
	//сделать всё, что нужно, пока тебе никто не мешает, и выключить обратно
}

var logger *log.Logger

func init() {
	logger = log.GetLogger("Room")
}

func (r *Room) TryJoin(p player.Player) (success bool) {
	//Здесь нам нужно под мьютексом проверить наличие свободных мест. Варианты:
	//1. 2 места свободно -> занимаем первое место
	//2. Свободно 1 место -> занимаем место, поднимаем флаг недоступности комнаты, начинаем игровой процесс
	logger.Infof("player %d entered Try Join", p.UID())
	r.mu.Lock()
	found := false

	if r.p1 == nil {
		r.p1 = p
		r.p1Status = StatusJoined
		logger.Infof("Player %d is now p1 in room", p.UID())
		found = true
	} else if r.p2 == nil {
		r.p2 = p
		r.p1Status = StatusJoined
		logger.Infof("Player %d is now p2 in room", p.UID())

		found = true
	}

	if r.p1 != nil && r.p2 != nil {
		logger.Infof("All players joined the game, p1: %d, p2: %d", r.p1.UID(), r.p2.UID())
		//TODO Prepare Match
		//TODO Then run buildEnv after PrepareMatch
		// In build Env составление и доставание даннных для вопросов
		go func() {
			//TODO handle possible panic
			r.prepareMatch()
		}()

	}

	r.mu.Unlock()
	return found
}

func (r *Room) notifyP1(msg messge.Message) (err error) {
	err = r.p1.Send(msg)
	if err != nil {
		logger.Error("Failed to send Message to P1", err)
	}
	return
}

func (r *Room) notifyP2(msg messge.Message) (err error) {
	err = r.p2.Send(msg)
	if err != nil {
		logger.Error("Failed to send Message to P2", err)
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

//Проверка Уместности сообщения ( на уровне комнаты)
func (r *Room) isSyncValid(wm MessageWrapper) (isValid bool) {
	r.mu.Lock()
	if wm.msg.Title == messge.Leave {
		isValid = true
		r.mu.Unlock()
		return
	}
	if wm.msg.Title == messge.ChangeOpponent || wm.msg.Title == messge.Quit || wm.msg.Title == messge.Continue {
		isValid = true
		r.mu.Unlock()
		return
	}

	if wm.player != r.active && (wm.msg.Title != messge.Ready) {
		logger.Error("isSync Player addr error")
		isValid = false
		r.mu.Unlock()
		return
	}
	if r.waitForSyncMsg != wm.msg.Title {
		logger.Error("isSync title error")
		isValid = false
		r.mu.Unlock()
		return
	}
	isValid = true
	r.mu.Unlock()
	return
}

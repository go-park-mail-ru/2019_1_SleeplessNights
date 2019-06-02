package factory

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type websocketPlayer struct {
	//С помощью этой структуры будем делать игрока по вебсокету на фабрике
	//Она уже реализует интерфейс Player, поэтому нам нужно будет просто сделатьинстанс этой структуры
	id   uint64
	uid  uint64
	in   chan message.Message
	conn *websocket.Conn
	mu   sync.Mutex
}

func (wsPlayer *websocketPlayer) StartListen() {
	//Метод, который заставляет структуру ждать сообщения
	//Если его не вызвать, то игрок не сможет сообщить серверу о своих действиях
	for {
		var msg message.Message

		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered ", r)
			}

		}()
		err := wsPlayer.conn.ReadJSON(&msg)
		if err != nil {
			//В случае получения ошибки нам нельзя прекращать слушать клиента, т.к. возможно, что
			//она была разовой и не критической
			//Нужно сформировать кастомое сообщение и отправить его серверу,
			// то-то типа "от кигрока пришло битое сообщение"

			if websocket.IsUnexpectedCloseError(err) {
				logger.Error("!!!!!!!!!!!!!!!!!!   PLAYER CLOSED CONNECTION   !!!!!!!!!!!!!!!!!!!!1")
				logger.Infof("Player %d closed the connection", wsPlayer.uid)
				//wsPlayer.in <- message.Message{Title: message.Leave}
				logger.Info("Before Close attempt in Start Listen Unxexpected Close")
				time.Sleep(time.Second)
				//wsPlayer.Close()
				logger.Info("After Close attempt in Start Listen Unxexpected Close")

				return
			}
		}
		logger.Info("Got from connection", msg)
		wsPlayer.in <- msg

	}
}

func (wsPlayer *websocketPlayer) Send(msg message.Message) (err error) {
	//Получаем наш месседж, который хотим отправить, и отправляем его в формате JSON
	wsPlayer.mu.Lock()
	err = wsPlayer.conn.WriteJSON(msg)
	if err != nil {
		return
	}
	wsPlayer.mu.Unlock()
	return nil
}

func (wsPlayer *websocketPlayer) Subscribe() chan message.Message {
	return wsPlayer.in
}

func (wsPlayer *websocketPlayer) ID() uint64 {
	return wsPlayer.id
}

func (wsPlayer *websocketPlayer) UID() uint64 {
	return wsPlayer.uid
}

func (wsPlayer *websocketPlayer) Close() {

	logger.Infof("Player UID  %d closed the connection", wsPlayer.uid)
	err := wsPlayer.conn.Close()
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		err := recover()
		if err != nil {
			logger.Error("(websocketPlayer) Close()", err)
		}
	}()
	wsPlayer.in <- message.Message{Title: message.Leave}
	close(wsPlayer.in)
}

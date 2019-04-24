package factory

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/messge"
	"github.com/gorilla/websocket"
)

type websocketPlayer struct {
	//TODO develop Close() method
	//С помощью этой структуры будем делать игрока по вебсокету на фабрике
	//Она уже реализует интерфейс Player, поэтому нам нужно будет просто сделатьинстанс этой структуры
	id        uint64
	uid       uint64
	in        chan messge.Message
	conn      *websocket.Conn
}

func (wsPlayer *websocketPlayer) StartListen() {
	//Метод, который заставляет структуру ждать сообщения
	//Если его не вызвать, то игрок не сможет сообщить серверу о своих действиях
	for {
		var msg messge.Message
		err := wsPlayer.conn.ReadJSON(&msg)
		if err != nil {
			//В случае получения ошибки нам нельзя прекращать слушать клиента, т.к. возможно, что
			//она была разовой и не критической
			//Нужно сформировать кастомое сообщение и отправить его серверу,
			// то-то типа "от кигрока пришло битое сообщение"
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logger.Infof("Player %d closed the connection", wsPlayer.uid)
				wsPlayer.in <- messge.Message{Title:messge.Leave}
				return
			}
		}
		logger.Info("Got from connection", msg)
		wsPlayer.in <- msg
	}
}

func (wsPlayer *websocketPlayer) Send(msg messge.Message) (err error) {
	//Получаем наш месседж, который хотим отправить, и отправляем его в формате JSON
	err = wsPlayer.conn.WriteJSON(msg)
	if err != nil {
		return
	}
	return nil
}

func (wsPlayer *websocketPlayer) Subscribe() chan messge.Message {
	return wsPlayer.in
}

func (wsPlayer *websocketPlayer) ID() uint64 {
	return wsPlayer.id
}

func (wsPlayer *websocketPlayer) UID() uint64 {
	return wsPlayer.uid
}

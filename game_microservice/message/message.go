package message

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Message")
}

const (
	//Набор констант, которые можно использовать в качестве значения поля Title для Message
	//На данном этапе трудно спрогнозировать полный набор таких заголовков,
	//поэтому значения приведены просто для примера и поменяются при реализации

	//ИСХОДЯЩИЕ
	StartGame    = "START_GAME"    // Оповещаем клиентов о том, что комната готова и они могут начать её отрисовывать
	YourTurn     = "YOUR_TURN"     // Оповещаем клиента о начале его хода
	OpponentTurn = "OPPONENT_TURN" // Оповещаемк клиента о том, что ходит его оппонент
	SelectedCell = "SELECTED_CELL" // выбранная для хода Клетка

	AvailableCells    = "AVAILABLE_CELLS"    // Оповещаем клиента о том, на какие клетки он может ходить; payload = []pair
	YourQuestion      = "QUESTION"           // Даём клиенту вопрос, связанный с клеткой; payload = question
	OpponentQuestion  = "OPPONENT_QUESTION"  // Оповещаем клиента о вопросе, на который отвечает его оппонент; payload = question
	OpponentAnswer    = "OPPONENT_ANSWER"    // Оповещаем клиента об ответе, который дал его оппонент; payload = int
	YourAnswer        = "YOUR_ANSWER"        // Оповещаем обоих клиентов о том, что какой ответ был выбран а какой был правильный
	Loss              = "LOSS"               // Оповещаем клиента о его поражении
	Win               = "WIN"                // Оповещаем клиента о его победе
	OpponentProfile   = "OPPONENT_PROFILE"   // Данные оппонента
	WannaPlayAgain    = "WANNA_PLAY_AGAIN"   // Даём клиенту выбор продолжить играть или нет
	OpponentLeaves    = "OPPONENT_QUITS"     // Оповещаем клиента о желании соперника продолжить
	OpponentContinues = "OPPONENT_CONTINUES" // Оповещаем клиента о желании выйти из игры

	CurrentState           = "CURRENT_STATE"           // Текущее состояние игры
	ThemesRequest          = "THEMES_REQUEST"          // Массив тем игрового поля
	QuestionsThemesRequest = "QUESTION_THEMES_REQUEST" // Массив id тем для вопросов

	AvailablePacks = "AVAILABLE_PACKS" // массив все возможных паков

	RoomSearching = "ROOM_SEARCHING" // Уведомление о начале поиска комната
	SelectedPack  = "SELECTED_PACK"

	//ВХОДЯЩИЕ
	//Входящие команды разделяются на синхронные (SYNC) и асинхронные (ASYNC)
	//Асинхронные команды всегда принимаются и добавляются в очередь входных сообщений комнаты
	//Синхронные команды могут быть отправлены только в установленном порядке (т.е. мы не можем обработать ANSWER раньше GO_TO)
	//TODO добавить валидатор входящих сообщений
	//TODO для этого добавить в комнату строковую переменную waitForSyncMsg, которая будет содержать заголовок следующей ожидаемой синхронной команды
	//TODO синхронные команды принимаются только от того игрока, чей сейчас ход
	//TODO получив синхронную команду, сравниваем её title с waitForSyncMsg, если совпадают, то добавляем в очередь и
	//TODO определяем для waitForSyncCmd новое значение (следующую ожидаемую команду). Если больше мы таких не ждём, то пустую строку
	Ready        = "READY"  // ASYNC Оповещаем сервер о том, что клиент подгрузился и можно стартовать таймер
	GoTo         = "GO_TO"  // SYNC Оповещаем клиента о клетке, которую выбрали для хода; payload = pair
	ClientAnswer = "ANSWER" // SYNC Оповещаем сервер о выбранном ответе на вопрос; payload = int
	Leave        = "LEAVE"  // ASYNC Оповещаем клиента о выходе из комнаты

	//Ответы игроков после того, как матч завершиться (ответы на запрос WannaPlayAgain)
	Quit           = "QUIT"            //  Оповещаем сервер о желании выйти из игры и в главное меню
	Continue       = "CONTINUE"        //  Оповещаем сервер о желании продолжить игру с тем же соперником
	ChangeOpponent = "CHANGE_OPPONENT" //  Оповещаем сервер о желании продолжить игру с другим соперником

	State = "STATE" // Запрос текущего состояния игры

	Themes          = "THEMES"          // Запрос  матрицы тем игрового поля
	QuestionsThemes = "QUESTION_THEMES" // Массив id тем для вопросов

	NotDesiredPack = "NOT_DESIRED_PACK" //  нежелательны pack
)

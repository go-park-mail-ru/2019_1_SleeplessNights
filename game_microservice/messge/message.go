package messge

import (
	"fmt"
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
	StartGame        = "START_GAME"         //Оповещаем клиентов о том, что комната готова и они могут начать её отрисовывать
	YourTurn         = "YOUR_TURN"          //Оповещаем клиента о начале его хода
	EnemyTurn        = "ENEMY_TURN"         //Оповещаемк клиента о том, что ходит его оппонент
	AvailableCells   = "AVAILABLE_CELLS"    //Оповещаем клиента о том, на какие клетки он может ходить; payload = []pair
	YourQuestion     = "QUESTION"           //Даём клиенту вопрос, связанный с клеткой; payload = question
	EnemyQuestion    = "ENEMY_QUESTION"     //Оповещаем клиента о вопросе, на который отвечает его оппонент; payload = question
	EnemyAnswer      = "ENEMY_ANSWER"       //Оповещаем клиента об ответе, который дал его оппонент; payload = int
	Correct          = "CORRECT"            //Оповещаем обоих клиентов о том, что ответ на вопрос верен
	Incorrect        = "INCORRECT"          //Оповещаем обоих клиентов о том, что ответ на вопрос неверен
	Loss             = "LOSS"               //Оповещаем клиента о его поражении
	Win              = "WIN"                //Оповещаем клиента о его победе
	OpponentProfile  = "OPPONENT_PROFILE"   // Данные оппонента
	WannaPlayAgain   = "WANNA_PLAY_AGAIN"   // Даём клиенту выбор продолжить играть или нет
	OpponentLeaves   = "OPPONENT_QUITS"     //Оповещаем клиента о желании соперника продолжить
	OpponentContines = "OPPONENT_CONTINUES" //Оповещаем клиента о желании выйти из игры

	CurrentState  = "CURRENT_STATE"  //Текущее состояние игры
	ThemesRequest = "THEMES_REQUEST" //матрица тем игрового поля

	//ВХОДЯЩИЕ
	//Входящие команды разделяются на синхронные (SYNC) и асинхронные (ASYNC)
	//Асинхронные команды всегда принимаются и добавляются в очередь входных сообщений комнаты
	//Синхронные команды могут быть отправлены только в установленном порядке (т.е. мы не можем обработать ANSWER раньше GO_TO)
	//TODO добавить валидатор входящих сообщений
	//TODO для этого добавить в комнату строковую переменную waitForSyncMsg, которая будет содержать заголовок следующей ожидаемой синхронной команды
	//TODO синхронные команды принимаются только от того игрока, чей сейчас ход
	//TODO получив синхронную команду, сравниваем её title с waitForSyncMsg, если совпадают, то добавляем в очередь и
	//TODO определяем для waitForSyncCmd новое значение (следующую ожидаемую команду). Если больше мы таких не ждём, то пустую строку
	Ready        = "READY"  //ASYNC Оповещаем сервер о том, что клиент подгрузился и можно стартовать таймер
	GoTo         = "GO_TO"  //SYNC Оповещаем клиента о клетке, которую выбрали для хода; payload = pair
	ClientAnswer = "ANSWER" //SYNC Оповещаем сервер о выбранном ответе на вопрос; payload = int
	Leave        = "LEAVE"  //ASYNC Оповещаем клиента о выходе из комнаты

	//Ответы игроков после того, как матч завершиться (ответы на запрос WannaPlayAgain)
	Quit           = "QUIT"            //  Оповещаем сервер о желании выйти из игры и в главное меню
	Continue       = "CONTINUE"        //  Оповещаем сервер о желании продолжить игру с тем же соперником
	ChangeOpponent = "CHANGE_OPPONENT" //  Оповещаем сервер о желании продолжить игру с другим соперником

	State          = "STATE"           //Запрос текущего состояния игры
	ThemesResponse = "THEMES_RESPONSE" //Запрос  матрицы тем игрового поля
)

type Message struct {
	//Формат пакета, средствами которых реализуется общение между клиентом и сервером
	//Самый простой вариант - JSON, и у нас нет причин от него отказываться
	//Можно было сделать с помощью интерфейса чтобы абстрагироваться от формата передаваемых данных,
	//но практического применения этому я не вижу

	//CommandName лишний, CommandName = Title

	Title   string      `json:"title"`
	Payload interface{} `json:"payload,omitempty"`
}

type Coordinates struct {
	//Achtung!!!!
	X int `json:"x"`
	Y int `json:"y"`
}
type ThemeArray struct {
	ThemesArray []string `json:"theme_array"`
}

//Request TryMove to a cell
type ThemePack struct {
	Id    uint   `json:"pack"`
	Theme string `json:"theme_name"`
}

//response from sever with question
type Question struct {
	Question string `json:"question"`
}
type GameState struct {
	State string `json:"state"`
}

//response from client with answer_id
type Answer struct {
	AnswerId int `json:"answer_id"`
}

//Response to players answer
type AnswerResult struct {
	AnswerResult  bool `json:"answer_id"`
	CorrectAnswer int  `json:"correct_answer"`
}

func (m *Message) IsValid() bool {
	fmt.Println(m.Payload)
	switch m.Title {
	case Ready:
		{
			return true
		}
	case GoTo:
		{

			st, ok := m.Payload.(map[string]interface{})
			if !ok {
				logger.Error("Message validator, Title=GO_TO, error:interface->Answer casting error")
				return false
			}
			if _, ok := st["x"]; !ok {
				return false
			}
			if _, ok := st["y"]; !ok {
				return false
			}
			return true

		}
	case ClientAnswer:
		{
			st, ok := m.Payload.(map[string]interface{})
			if !ok {
				logger.Error("Message validator, Title=ClientAnswer, error:interface->Answer casting error")
				return false
			}
			if _, ok := st["answer_id"]; !ok {
				return false
			}

			return true
		}
	case Leave:
		{
			return true
		}

	case Continue:
		{
			return true
		}
	case ChangeOpponent:
		{
			return true
		}
	case Quit:
		{
			return true
		}
	case State:
		{
			return true
		}
	case ThemesRequest:
		{
			return true
		}
	default:
		return false
	}
}

package messge

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/questions"

const (
	//Набор констант, которые можно использовать в качестве значения поля Title для Message
	//На данном этапе трудно спрогнозировать полный набор таких заголовков,
	//поэтому значения приведены просто для примера и поменяются при реализации
	command = "COMMAND"
	info    = "INFO"

	CommandMove   = "CommandMove"
	CommandAnswer = "CommandAnswer"

	//TODO разработать API
	/***/
)

//PAYLOAD FORMATS

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

//Request TryMove to a cell

type MoveRequest struct {
	PlayerId        uint64      `json:"player_id"`
	CurrentPosition Coordinates `json:"curr_pos"`
	DesiredPosition Coordinates `json:"desired_pos"`
}

//response from sever with question
type Question struct {
	PlayerId uint64             `json:"player_id"`
	Question questions.Question `json:"question"`
}

//response from client with answer_id
type Answer struct {
	PlayerId uint64 `json:"player_id"`
	AnswerId int    `json:"answer_id"`
}

//Response to players answer
type AnswerResult struct {
	PlayerId     uint64 `json:"player_id"`
	AnswerResult bool   `json:"answer_id"`
}

type Message struct {
	//Формат пакета, средствами которых реализуется общение между клиентом и сервером
	//Самый простой вариант - JSON, и у нас нет причин от него отказываться
	//Можно было сделать с помощью интерфейса чтобы абстрагироваться от формата передаваемых данных,
	//но практического применения этому я не вижу
	Title       string      `json:"title"`
	CommandName string      `json:"command_name, omitempty"`
	Payload     interface{} `json:"payload,omitempty"`
}

package message
//go:generate $GOPATH/bin/easyjson structs.go
//easyjson:json
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

type GameState struct {
	State string `json:"state"`
}

//response from client with answer_id
type Answer struct {
	AnswerId int `json:"answer_id"`
}

//Response to players answer
type AnswerResult struct {
	GivenAnswer   int `json:"given_answer"`
	CorrectAnswer int `json:"correct_answer"`
}

type PackID struct {
	PackId int64 `json:"pack_id"`
}

package event

type eventType int

const (
	Move = iota
	WinPrize
	Lose
	Info
	NoMoves
	Warning
	Error
)

//Event sent by "gameField" as answer to "Room's" call of gf methods
type Event struct {
	Etype eventType
	Edata interface{}
}

type Pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Coordinates []Pair

type Question struct {
	Question string `json:"question"`
}

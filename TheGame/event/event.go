package event

type eventType int

const (
	Move    = iota
	Win     = iota
	Lose    = iota
	Info    = iota
	NoMoves = iota
)

//Event sent by "gameField" as answer to "Room's" call of gf methods
type Event struct {
	Etype eventType

	Edata interface{}
}

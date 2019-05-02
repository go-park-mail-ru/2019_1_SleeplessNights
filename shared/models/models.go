package models

type Question struct {
	ID      uint64
	Answers []string
	Correct int
	Text    string
	PackID  uint
}

type Pack struct {
	ID    uint64
	Theme string
}

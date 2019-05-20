package database

type Question struct {
	ID      uint64   `json:"-"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"`
	Text    string   `json:"text"`
	PackID  uint     `json:"pack_id"`
	JSON    []byte   `json:"-"`
}

type Pack struct {
	ID       uint64 `json:"id"`
	IconPath string `json:"iconPath"`
	Theme    string `json:"name"`
}

type QuestionForFrontend struct {
	Text    string   `json:"text"`
	Answers []string `json:"answers"`
	PackID  uint     `json:"pack_id"`
}

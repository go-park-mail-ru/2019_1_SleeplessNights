package database
//go:generate $GOPATH/bin/easyjson models.go
type Question struct {
	ID      uint64   `json:"-"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"`
	Text    string   `json:"text"`
	PackID  uint     `json:"pack_id"`
	JSON    []byte   `json:"-"`
}

type Pack struct {
	ID uint64 `json:"id"`
	Theme    string `json:"name"`
	IconPath string `json:"iconPath"`
}
//easyjson:json
type QuestionForFrontend struct {
	Text    string   `json:"text"`
	Answers []string `json:"answers"`
	PackID  uint     `json:"pack_id"`
}

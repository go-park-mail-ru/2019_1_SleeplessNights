package database

import "fmt"

type Question struct {
	ID      uint64   `json:"-"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"`
	Text    string   `json:"text"`
	PackID  uint     `json:"pack_id"`
}

type Pack struct {
	ID       uint64 `json:"id"`
	IconPath string `json:"iconPath"`
	Theme    string `json:"theme"`
}

type QuestionForFrontend struct {
	Text    string   `json:"text"`
	Answers []string `json:"answers"`
	PackID  uint     `json:"pack_id"`
}

func (q *Question) ToJson() (json string) {
	json = "{"

	json += `"text:"` + q.Text + ","
	json += `"answers":[`
	for _, answer := range q.Answers {
		json += answer + ","
	}
	json += "]"
	json += "}"
	fmt.Println(json)
	return
}

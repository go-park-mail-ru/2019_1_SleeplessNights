package models

import "fmt"

type Question struct {
	ID      uint64   `json:"-"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"`
	Text    string   `json:"text"`
	PackID  uint     `json:"pack_id"`
}

type Pack struct {
	ID    uint64 `json:"id"`
	Theme string `json:"theme"`
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

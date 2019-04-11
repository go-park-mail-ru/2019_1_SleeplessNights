package questions

type Question struct {
	QuestionJson    string `json:"question"`
	CorrectAnswerId int    `json:"-"`
	//Здесь ситуация аналогичная с гейм филдом
	//Сама структура из себя представляет JSON, который мы прямо так и будем хранить в БД,
	//и номер правильного ответа
	//Ответственность:
	//* Давать по запросу JSON вопроса (варианты ответа уже лежат внутри)
	//* Давать по запросу номер правильного ответа
	//TODO develop
}

func (q *Question) GetQuestion() (question string) {
	return q.QuestionJson
}

func (q *Question) CheckAnswer(answerId int) (result bool) {
	return q.CorrectAnswerId == answerId
}

func (q *Question) GetAnswerId() (answer_id int) {
	return q.CorrectAnswerId
}

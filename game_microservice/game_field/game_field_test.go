package game_field

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"math/rand"
	"sync/atomic"
	"testing"
)

func TestGameField_Build(t *testing.T) {
	gf := GameField{}
	var idSource uint64
	idSource = 0
	getRandQuestion := func(idSource *uint64) database.Question {
		atomic.AddUint64(idSource, 1)
		return database.Question{
			PackID:          *idSource,
			QuestionJ:    "{}",
			CorrectAnswerId: rand.Int() % 4,
		}
	}

	var qArray [QuestionsNum]database.Question
	for i := range qArray {
		qArray[i] = getRandQuestion(&idSource)
	}

	gf.Build(qArray)

	result := gf.field

	getAnswer := func(id uint64) int {
		for i := 0; i < fieldSize; i++ {
			for j := 0; j < fieldSize; j++ {
				if !result[i][j].isAvailable {
					t.Errorf("Cell (%d, %d) is not available by default", i, j)
				}
				if !isPrizePosition(i, j) && result[i][j].question.PackID == id {
					return result[i][j].question.CorrectAnswerId
				}
			}
		}
		return -1
	}

	for _, q := range qArray {
		if getAnswer(q.PackID) != q.CorrectAnswerId {
			t.Error("Build method violates questions")
			return
		}
	}
}

func TestGameField_CheckAnswer(t *testing.T) {

}

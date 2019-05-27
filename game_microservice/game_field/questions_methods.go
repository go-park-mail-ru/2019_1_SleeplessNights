package game_field

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"math/rand"
	"time"
)

func Shuffle(questions *[]database.Question) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(*questions); n > 0; n-- {
		randIndex := r.Intn(n)
		(*questions)[n-1], (*questions)[randIndex] = (*questions)[randIndex], (*questions)[n-1]
	}

}

func (gf *GameField) GetQuestionsThemes() (packArray []uint) {
	for i := 0; i < fieldSize; i++ {
		for j := 0; j < fieldSize; j++ {
			if gf.field[i][j].question != nil {
				fmt.Println((gf.field[i][j]).question)
				packArray = append(packArray, ((gf.field[i][j]).question.PackID))
			}
		}
	}
	return
}

func (gf *GameField) GetQuestionByCell(x, y int) (question *database.Question) {
	logger.Infof("GetQuestionByCell x:%d,y:%d ", x, y)
	question = gf.field[x][y].question
	return
}

func (gf *GameField) CheckAnswer(answerIdx int) bool {
	if !gf.validateAnswerId(answerIdx) {
		return false
	}
	if gf.regQuestion.Correct == answerIdx {
		(gf.field[gf.regY][gf.regX]).isAvailable = false
		(gf.field[gf.regY][gf.regX]).answerResult = 1
		return true
	}
	(gf.field[gf.regY][gf.regX]).isAvailable = false
	(gf.field[gf.regY][gf.regX]).answerResult = -1
	return false
}

func (gf *GameField) GetRegisterQuestion() database.Question {
	return gf.regQuestion
}

func (gf *GameField) GetPacksSlice() (packs *[]database.Pack) {
	return &gf.Themes
}

func (gf *GameField) SetPacksSlice(packs []database.Pack) {
	copy(gf.Themes, packs)
}

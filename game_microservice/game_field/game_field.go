package game_field

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"strconv"
)

var logger *log.Logger

const fieldSize = 8

var (
	QuestionsNum = config.GetInt("game_ms.pkg.game_field.questions_num")
	TurnDuration = config.GetInt("game_field.pkg.game_field.turn_duration")
)

var prizePos []Pair

//В начале игры игроков не существует никаких позиций, они находятся как бы вне поля
func init() {
	logger = log.GetLogger("GameMS")
	logger.SetLogLevel(logrus.Level(config.GetInt("game_ms.log_level")))
	prizePosRaw := config.GetMapStringToInterface("game_ms.pkg.game_field.prize_pos")
	for i := 0; i < len(prizePosRaw); i++ {
		prize, ok := prizePosRaw[strconv.Itoa(i)].(map[string]interface{})
		if !ok {
			logger.Fatal("Unable to convert prize to map[string]interface{}")
		}
		x, err := strconv.Atoi(prize["x"].(string))
		if err != nil {
			logger.Fatal("Can't convert prizePos.x to int")
		}
		y, err := strconv.Atoi(prize["y"].(string))
		if err != nil {
			logger.Fatal("Can't convert prizePos.y to int")
		}
		prizePos = append(prizePos, Pair{X: x, Y: y})
	}
}

type GameField struct {
	Themes []database.Pack
	field  [fieldSize][fieldSize]gameCell
	p1     gfPlayer
	p2     gfPlayer
	//Out   []event.Event
	regX        int
	regY        int
	regQuestion database.Question

	//Тут уровень абстракции уже достаточно постой для понимания, поэтому оставляю реализацию на ваше усмотрение
	//По ответственности, если навскидку, игровое поле должно:
	//* Знать расположение всех вопросов
	//* Знать координаты всех игроков (при этом самих игроков знать не должно)
	//* Знать состояние каждой клетки поля: доступна/недоступна для хода
	//* Уметь проверять, что кто-то из игроков уже достиг центра и победил
	//* Уметь проверять, что кто-то из игроков оказался недостижим от центра поля и проиграл (задачка со звёздочкой на алгоритм Дейкстры)
	//TODO develop
}

//Cell coordinates are taken from gamefield register

func (gf *GameField) ResetPlayersPositions() {
	gf.p1.pos = nil
	gf.p2.pos = nil
}

func (gf *GameField) SetPlayerFirstPosition() {
	gf.p1.pos = &Pair{1, 1}
	return
}

func (gf *GameField) GetCurrentState() string {
	fieldState := fmt.Sprintln("\n x0 x1 x2 x3 x4 x5 x6 x7\n __ __ __ __ __ __ __ __")
	for i := 0; i < fieldSize; i++ {
		for j := 0; j < fieldSize; j++ {
			if gf.p1.pos != nil {
				if (*gf.p1.pos) == (Pair{j, i}) {
					fieldState = fieldState + "|P1"
					continue
				}
			}
			if gf.p2.pos != nil {
				if (*gf.p2.pos) == (Pair{j, i}) {
					fieldState = fieldState + "|P2"
					continue
				}
			}
			if isPrizePosition(j, i) {
				fieldState = fieldState + "|Pr"
				continue
			}
			if gf.field[i][j].answerResult == 1 {
				fieldState = fieldState + "|+_"
			}
			if gf.field[i][j].answerResult == -1 {
				fieldState = fieldState + "|-_"
			}
			if gf.field[i][j].answerResult == 0 {
				fieldState = fieldState + "|__"
			}
			if j == 7 {
				fieldState = fieldState + fmt.Sprintln("|  y", i)
			}
			if i == 7 && j == 7 {
				fieldState = fieldState + fmt.Sprintln()
			}
		}
	}
	p1State := ""
	p2State := ""
	if gf.p1.pos != nil && gf.p2.pos != nil {
		p1State = fmt.Sprintf("\n\n player1 %d, {X:%d, Y:%d},answers: +:%d -:%d \n", gf.p1.id, gf.p1.pos.X, gf.p1.pos.Y, gf.p1.rightAnswers, gf.p1.falseAnswers)
		p2State = fmt.Sprintf("\n\n player1 %d, {X:%d, Y:%d},answers: +:%d -:%d \n ", gf.p2.id, gf.p2.pos.X, gf.p2.pos.Y, gf.p1.rightAnswers, gf.p1.falseAnswers)

	}
	return fieldState + p1State + p2State
}

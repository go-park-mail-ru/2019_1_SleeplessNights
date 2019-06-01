package user_manager

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

const (
	nodataFound     = "P0002"
	uniqueViolation = "23505"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("User")
	logger.SetLogLevel(logrus.Level(config.GetInt("user_ms.log_level")))
}

const defaultLeaderBoardUpdateInterval = 1 * time.Second

var LeaderBoardLen = uint64(config.GetInt("user_ms.pkg.user_manager.board_len"))

var profiles []*services.Profile

var user *userManager

type userManager struct {
	secret []byte
}

func init() {
	//secretFile, err := os.Open(os.Getenv("BASEPATH") + "/secret")
	secretFile, err := os.Open( "/Users/mac/Desktop/back-end/2019_1_SleeplessNights/secret")

	defer func() {
		err := secretFile.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	if err != nil {
		logger.Fatal(err)
	}

	var secret []byte
	_, err = fmt.Fscanln(secretFile, &secret)
	if err != nil {
		return
	}

	user = &userManager{
		secret: secret,
	}
}

func GetInstance() *userManager {
	return user
}


func UpdateLeaderBoard()  {
	leaderBoardUpdateInterval, err := time.ParseDuration(config.GetString("user_ms.pkg.user_manager.leaderboard_update_interval"))
	if err != nil {
		leaderBoardUpdateInterval = defaultLeaderBoardUpdateInterval
	}

	for {
		wg := sync.WaitGroup{}
		wg.Add(1)
		time.AfterFunc(leaderBoardUpdateInterval, func() {
			profiles, err = database.GetInstance().GetUsers(LeaderBoardLen)
			if err != nil {
				logger.Errorf("Failed to get users: %v", err.Error())
				return
			}
			logger.Info("Updated successful profiles")
			leaderBoardUpdateInterval, err = time.ParseDuration(config.GetString("user_ms.pkg.user_manager.leaderboard_update_interval"))
			if err != nil {
				leaderBoardUpdateInterval = defaultLeaderBoardUpdateInterval
			}
			wg.Done()
		})
		wg.Wait()
	}
}

func handlerError(pgError pgx.PgError) (err error) {
	switch pgError.Code {
	case uniqueViolation:
		err = errors.DataBaseUniqueViolation
	case nodataFound:
		err = errors.DataBaseNoDataFound
	default:
		err = pgError
	}
	return
}

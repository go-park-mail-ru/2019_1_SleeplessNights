package user_manager

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"os"
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

var user *userManager

type userManager struct {
	secret []byte
}

func init() {
	secretFile, err := os.Open(os.Getenv("BASEPATH") + "/secret") //TODO а тут нужно уберать basepath????
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

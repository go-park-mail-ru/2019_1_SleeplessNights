package user_manager

import (
	"fmt"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	sessionLifeLen = 4 * time.Hour
	NoTokenOwner   = "error: There are no token's owner in database"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Auth")
	logger.SetLogLevel(logrus.TraceLevel)
}

var user *userManager

type userManager struct {
	secret []byte
}

func init() {
	secretFile, err := os.Open(os.Getenv("BASEPATH") + "/secret")
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

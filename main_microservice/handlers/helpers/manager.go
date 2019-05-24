package helpers

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"regexp"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Validator")
	logger.SetLogLevel(logrus.TraceLevel)
}

var (
	emailReg *regexp.Regexp
	nicknameReg *regexp.Regexp
)

func init(){
	var err error
	emailReg, err = regexp.Compile("^[a-z0-9._%+-]+@[a-z0-9-]+.+.[a-z]{2,4}$")
	if err != nil {
		logger.Errorf("Failed to match: %v", err.Error())
		return
	}

	nicknameReg, err = regexp.Compile("^[a-zA-Z0-9-_]*$")
	if err != nil {
		logger.Errorf("Failed to match: %v", err.Error())
		return
	}
}
package errors

import (
	"errors"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"regexp"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Errors")
	logger.SetLogLevel(logrus.TraceLevel)
}

var (
	DataBaseUniqueViolation     = errors.New("ERROR: unique violation exception in database")
	DataBaseNoDataFound         = errors.New("ERROR: no data found exception in database")
	DataBaseForeignKeyViolation = errors.New("ERROR: foreign key violation exception in database")
	AuthWrongPassword           = errors.New("ERROR: authentication failed, because of wong password")
)

var (
	DataBaseUniqueViolationReg     *regexp.Regexp
	DataBaseNoDataFoundReg         *regexp.Regexp
	DataBaseForeignKeyViolationReg *regexp.Regexp
)

func init() {
	var err error
	DataBaseUniqueViolationReg, err = regexp.Compile(DataBaseUniqueViolation.Error())
	if err != nil {
		logger.Errorf("Failed to match: %v", err.Error())
		return
	}

	DataBaseNoDataFoundReg, err = regexp.Compile(DataBaseNoDataFound.Error())
	if err != nil {
		logger.Errorf("Failed to match: %v", err.Error())
		return
	}

	DataBaseForeignKeyViolationReg, err = regexp.Compile(DataBaseForeignKeyViolation.Error())
	if err != nil {
		logger.Errorf("Failed to match: %v", err.Error())
		return
	}
}

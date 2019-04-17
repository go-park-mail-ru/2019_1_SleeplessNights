package main_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/xlab/closer"
	"testing"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Tests")
}

func TestMain(m *testing.M) {
	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		logger.Error(err.Error())
	}
	defer closer.Close()
	m.Run()
}

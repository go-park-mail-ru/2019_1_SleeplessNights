package logger

import (
	"log"
	"os"
)

var (
	Info    *log.Logger
	Error   *log.Logger
	Trace   *log.Logger
	Debug   *log.Logger
	Warning *log.Logger
	Fatal   *log.Logger
)

var (
	file = os.Stderr
)

func init() {
	Info = log.New(file,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Error = log.New(file,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Trace = log.New(file,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Debug = log.New(file,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Warning = log.New(file,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Fatal = log.New(file,
		"FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}

/*В любой пакет нужно испортировать ""github.com/go-park-mail-ru/2019_1_SleeplessNights/log""
log.<log>.Println("commit")*/

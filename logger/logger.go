package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

const fatalExitSuffix = "Exiting process with non zero code..."

type Logger struct {
	log    *logrus.Logger
	source string
}

func (logger *Logger)makeLog() *logrus.Entry {
	return logger.log.WithField("SOURCE", logger.source)
}

func (logger *Logger)InfoWithFields(fields logrus.Fields, msg ...interface{})  {
	logger.makeLog().WithFields(fields).Infoln(msg)
}

func (logger *Logger)InfofWithFields(fields logrus.Fields, format string, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Infof(format, msg)
}

func (logger *Logger)DebugWithFields(fields logrus.Fields, msg ...interface{})  {
	logger.makeLog().WithFields(fields).Debugln(msg)
}

func (logger *Logger)DebugfWithFields(fields logrus.Fields, format string, msg ...interface{})  {
	logger.makeLog().WithFields(fields).Debugf(format, msg)
}

func (logger *Logger)TraceWithFields(fields logrus.Fields, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Traceln(msg)
}

func (logger *Logger)TracefWithFields(fields logrus.Fields, format string, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Tracef(format, msg)
}

func (logger *Logger)WarningWithFields(fields logrus.Fields, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Warningln(msg)
}

func (logger *Logger)WarningfWithFields(fields logrus.Fields, format string, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Warningf(format, msg)
}

func (logger *Logger)ErrorWithFields(fields logrus.Fields, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Errorln(msg)
}

func (logger *Logger)ErorfWithFields(fields logrus.Fields, format string, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Errorf(format, msg)
}

func (logger *Logger)FatalWithFields(fields logrus.Fields, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Fatalln(msg, fatalExitSuffix)
}

func (logger *Logger)FatalfWithFields(fields logrus.Fields, format string, msg ...interface{}) {
	logger.makeLog().WithFields(fields).Fatalf(format, msg, fatalExitSuffix)
}

func (logger *Logger)Info(msg ...interface{})  {
	logger.makeLog().Infoln(msg)
}

func (logger *Logger)Infof(format string, msg ...interface{})  {
	logger.makeLog().Infof(format, msg)
}

func (logger *Logger)Debug(msg ...interface{})  {
	logger.makeLog().Debugln(msg)
}

func (logger *Logger)Debugf(format string, msg ...interface{})  {
	logger.makeLog().Debugf(format, msg)
}

func (logger *Logger)Trace(msg ...interface{})  {
	logger.makeLog().Traceln(msg)
}

func (logger *Logger)Tracef(format string, msg ...interface{})  {
	logger.makeLog().Tracef(format, msg)
}

func (logger *Logger)Warning(msg ...interface{})  {
	logger.makeLog().Warningln(msg)
}

func (logger *Logger)Warningf(format string, msg ...interface{})  {
	logger.makeLog().Warningf(format, msg)
}

func (logger *Logger)Error(msg ...interface{})  {
	logger.makeLog().Errorln(msg)
}

func (logger *Logger)Erorf(format string, msg ...interface{})  {
	logger.makeLog().Errorf(format, msg)
}

func (logger *Logger)Fatal(msg ...interface{})  {
	logger.makeLog().Fatalln(msg, fatalExitSuffix)
}

func (logger *Logger)Fatalf(format string, msg ...interface{})  {
	logger.makeLog().Fatalf(format, msg, fatalExitSuffix)
}

func GetLogger(source string) *Logger {
	logger := Logger{
		log:    logrus.New(),
		source: source,
	}

	logger.log.Formatter = new(logrus.TextFormatter)
	logger.log.Formatter.(*logrus.TextFormatter).TimestampFormat = "02-01-2006 15:04:05"
	logger.log.Formatter.(*logrus.TextFormatter).FullTimestamp = true
	logger.log.Out = os.Stderr
	return &logger
}

func (logger *Logger)SetLogLevel(level logrus.Level) {
	logger.log.Level = level
}


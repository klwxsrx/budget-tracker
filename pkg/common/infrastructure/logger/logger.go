package logger

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"

	"github.com/sirupsen/logrus"

	"os"
)

type logrusLogger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	Debug(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Info(args ...interface{})
	Fatal(args ...interface{})
}

type impl struct {
	logger logrusLogger
}

func (i *impl) With(fields log.Fields) log.Logger {
	return fromLogrus(i.logger.WithFields(logrus.Fields(fields)))
}

func (i *impl) WithError(err error) log.Logger {
	return fromLogrus(i.logger.WithError(err))
}

func (i *impl) Debug(args ...interface{}) {
	i.logger.Debug(args)
}

func (i *impl) Error(args ...interface{}) {
	i.logger.Error(args)
}

func (i *impl) Warn(args ...interface{}) {
	i.logger.Warning(args)
}

func (i *impl) Info(args ...interface{}) {
	i.logger.Info(args)
}

func (i *impl) Fatal(args ...interface{}) {
	i.logger.Fatal(args)
}

func fromLogrus(l logrusLogger) log.Logger {
	return &impl{l}
}

func New() log.Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetLevel(logrus.InfoLevel)

	return fromLogrus(logrusLogger)
}

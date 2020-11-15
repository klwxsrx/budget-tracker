package logger

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	external "github.com/sirupsen/logrus"
)

type externalLogger interface {
	WithFields(fields external.Fields) *external.Entry
	Error(args ...interface{})
	Warning(args ...interface{})
	Info(args ...interface{})
	Fatal(args ...interface{})
}

type impl struct {
	logger externalLogger
}

func (i *impl) With(fields logger.Fields) logger.Logger {
	return New(i.logger.WithFields(external.Fields(fields)))
}

func (i *impl) Error(args ...interface{}) {
	i.logger.Error(args)
}

func (i *impl) Warning(args ...interface{}) {
	i.logger.Warning(args)
}

func (i *impl) Info(args ...interface{}) {
	i.logger.Info(args)
}

func (i *impl) Fatal(args ...interface{}) {
	i.logger.Fatal(args)
}

func New(l externalLogger) logger.Logger {
	return &impl{l}
}

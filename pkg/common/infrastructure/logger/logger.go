package logger

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	external "github.com/sirupsen/logrus"
)

type externalLogger interface {
	WithFields(fields external.Fields) *external.Entry
	Info(args ...interface{})
}

type impl struct {
	logger externalLogger
}

func (i *impl) With(fields logger.Fields) logger.Logger {
	return New(i.logger.WithFields(external.Fields(fields)))
}

func (i *impl) Info(args ...interface{}) {
	i.logger.Info(args)
}

func New(l externalLogger) logger.Logger {
	return &impl{l}
}

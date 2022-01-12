package logger

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"

	external "github.com/sirupsen/logrus"

	"os"
)

type externalLogger interface {
	WithFields(fields external.Fields) *external.Entry
	WithError(err error) *external.Entry
	Debug(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Info(args ...interface{})
	Fatal(args ...interface{})
}

type impl struct {
	logger externalLogger
}

func (i *impl) With(fields logger.Fields) logger.Logger {
	return fromExternal(i.logger.WithFields(external.Fields(fields)))
}

func (i *impl) WithError(err error) logger.Logger {
	return fromExternal(i.logger.WithError(err))
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

func fromExternal(l externalLogger) logger.Logger {
	return &impl{l}
}

func New() logger.Logger {
	l := external.New()
	l.SetFormatter(&external.JSONFormatter{
		DisableHTMLEscape: true,
	})
	l.SetOutput(os.Stdout)
	l.SetLevel(external.InfoLevel)

	return fromExternal(l)
}

package pulsar

import (
	"fmt"
	pulsarLog "github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
)

type loggerAdapter struct {
	logger logger.Logger
}

func (l *loggerAdapter) SubLogger(fields pulsarLog.Fields) pulsarLog.Logger {
	return &loggerAdapter{l.logger.With(logger.Fields(fields))}
}

func (l *loggerAdapter) WithFields(fields pulsarLog.Fields) pulsarLog.Entry {
	return &loggerAdapter{l.logger.With(logger.Fields(fields))}
}

func (l *loggerAdapter) WithField(name string, value interface{}) pulsarLog.Entry {
	return &loggerAdapter{l.logger.With(logger.Fields{name: value})}
}

func (l *loggerAdapter) WithError(err error) pulsarLog.Entry {
	return &loggerAdapter{l.logger.WithError(err)}
}

func (l *loggerAdapter) Debug(args ...interface{}) {
	l.logger.Debug(args)
}

func (l *loggerAdapter) Info(args ...interface{}) {
	l.logger.Info(args)
}

func (l *loggerAdapter) Warn(args ...interface{}) {
	l.logger.Warn(args)
}

func (l *loggerAdapter) Error(args ...interface{}) {
	l.logger.Error(args)
}

func (l *loggerAdapter) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args))
}

func (l *loggerAdapter) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args))
}

func (l *loggerAdapter) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args))
}

func (l *loggerAdapter) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args))
}

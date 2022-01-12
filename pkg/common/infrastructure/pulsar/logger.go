package pulsar

import (
	"fmt"

	pulsarlog "github.com/apache/pulsar-client-go/pulsar/log"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
)

type loggerAdapter struct {
	logger log.Logger
}

func (l *loggerAdapter) SubLogger(fields pulsarlog.Fields) pulsarlog.Logger {
	return &loggerAdapter{l.logger.With(log.Fields(fields))}
}

func (l *loggerAdapter) WithFields(fields pulsarlog.Fields) pulsarlog.Entry {
	return &loggerAdapter{l.logger.With(log.Fields(fields))}
}

func (l *loggerAdapter) WithField(name string, value interface{}) pulsarlog.Entry {
	return &loggerAdapter{l.logger.With(log.Fields{name: value})}
}

func (l *loggerAdapter) WithError(err error) pulsarlog.Entry {
	return &loggerAdapter{l.logger.WithError(err)}
}

func (l *loggerAdapter) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *loggerAdapter) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *loggerAdapter) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *loggerAdapter) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *loggerAdapter) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *loggerAdapter) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *loggerAdapter) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *loggerAdapter) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

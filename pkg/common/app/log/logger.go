package log

type Fields map[string]interface{}

type Logger interface {
	With(fields Fields) Logger
	WithError(err error) Logger
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Fatal(args ...interface{})
}

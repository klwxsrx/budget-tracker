package logger

type Fields map[string]interface{}

type Logger interface {
	With(fields Fields) Logger
	Error(args ...interface{})
	Info(args ...interface{})
}

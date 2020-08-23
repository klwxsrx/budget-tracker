package logger

type Fields map[string]interface{}

type Logger interface {
	With(fields Fields) Logger
	Info(args ...interface{})
}

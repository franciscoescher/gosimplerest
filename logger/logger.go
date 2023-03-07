package logger

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type BlankLogger struct{}

func (b *BlankLogger) Info(args ...interface{}) {}

func (b *BlankLogger) Error(args ...interface{}) {}

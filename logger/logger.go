package logger

type Logger interface {
	Info(msg string, kv ...interface{})
	Warn(msg string, kv ...interface{})
	Error(msg string, kv ...interface{})
	Sync()
}

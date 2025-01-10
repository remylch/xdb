package log

import "time"

type LogLevel string

const (
	INFO  LogLevel = "info"
	ERROR          = "error"
)

type Log struct {
	timestamp time.Time
	level     LogLevel
	message   string
}

type Logger interface {
	send() error
}

type LogManager struct {
	logDir string
}

func NewLogManager(logDir string) Logger {
	return LogManager{
		logDir: logDir,
	}
}

func (l LogManager) send() error {
	return nil
}

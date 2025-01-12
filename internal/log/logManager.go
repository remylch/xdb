package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"xdb/internal/shared"

	"github.com/google/uuid"
)

type LogLevel string

const (
	INFO  LogLevel = "info"
	ERROR          = "error"
)

const (
	MaxLogFileSize = 1024 * 1024 * 10
	logFile        = "xdb.log"
)

type Log struct {
	timestamp     time.Time
	level         LogLevel
	message       string
	correlationId uuid.UUID
}

type Logger interface {
	Log(level LogLevel, messages ...string)
	Dir() string
}

type LogManager struct {
	logDir string
	mu     sync.RWMutex
}

func NewDefaultLogger(logDir string) Logger {
	if !shared.DirExists(logDir) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			panic(err)
		}

		if err := os.WriteFile(getLogFilePath(logDir), nil, os.ModePerm); err != nil {
			panic(err)
		}
	}

	return &LogManager{
		logDir: logDir,
		mu:     sync.RWMutex{},
	}
}

func (l *LogManager) Log(level LogLevel, messages ...string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	log := Log{
		timestamp: time.Now(),
		level:     level,
		message:   strings.Join(messages, " "),
	}

	msg := fmt.Sprintf("%s [%s] %s\n", log.timestamp.Format(time.RFC3339), log.level, log.message)

	fmt.Println(msg)

	if err := l.writeLog(msg); err != nil {
		fmt.Errorf("Error writing log to file : %v \n", err)
	}

}

func (l *LogManager) Dir() string {
	return l.logDir
}

func (l *LogManager) writeLog(log string) error {
	file, err := os.OpenFile(getLogFilePath(l.logDir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.WriteString(log); err != nil {
		return err
	}

	return nil
}

func getLogFilePath(logDir string) string {
	return logDir + "/" + logFile
}

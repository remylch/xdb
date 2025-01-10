package log

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
	"xdb/internal/shared"
)

type LogLevel string

const (
	INFO  LogLevel = "info"
	ERROR          = "error"
)

const (
	MaxLogFileSize = 1024 * 1024 * 10
	logFile = "xdb.log"
)

type Log struct {
	timestamp     time.Time
	level         LogLevel
	message       string
	correlationId uuid.UUID
}

type Logger interface {
	log(level LogLevel, message string)
}

type LogManager struct {
	logDir string
	mu sync.RWMutex
}

func NewLogManager(logDir string) Logger {
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
		mu: sync.RWMutex{},
	}
}

func (l *LogManager) log(level LogLevel, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	log := Log{
		timestamp: time.Now(),
		level:     level,
		message:   message,
	}

	if err := l.writeLog(log); err != nil {
		fmt.Errorf("Error writing log to file : %v \n", err)
	}

}

func (l *LogManager) writeLog(log Log) error {
	file, err := os.OpenFile(getLogFilePath(l.logDir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.WriteString(fmt.Sprintf("%s [%s] %s\n", log.timestamp.Format(time.RFC3339), log.level, log.message)); err != nil {
		return err
	}

	return nil
}

func getLogFilePath(logDir string) string {
	return logDir + "/" + logFile
}

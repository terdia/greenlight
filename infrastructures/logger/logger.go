package logger

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) PrintInfo(message string, context map[string]string) {
	l.print(LevelInfo, message, context)
}

func (l *Logger) PrintError(err error, context map[string]string) {
	l.print(LevelError, err.Error(), context)
}

func (l *Logger) PrintFatal(err error, context map[string]string) {
	l.print(LevelFatal, err.Error(), context)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, context map[string]string) (int, error) {

	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level   string            `json:"level"`
		Time    string            `json:"time"`
		Message string            `json:"message"`
		Context map[string]string `json:"context,omitempty"`
		Trace   string            `json:"trace,omitempty"`
	}{
		Level:   level.String(),
		Time:    time.Now().UTC().Format(time.RFC3339),
		Message: message,
		Context: context,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	// Lock the mutex so that no two writes to the output destination can happen
	// concurrently. If we don't do this, it's possible that the text for two or more
	// log entries will be intermingled in the output.
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}

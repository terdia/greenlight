package logger

import (
	"log"
	"os"
)

type logs struct {
	*log.Logger
}

func NewLogger() logs {
	return logs{Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime)}
}

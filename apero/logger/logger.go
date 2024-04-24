package logger

import (
	"log"
	"os"
)

func NewLogger(scope string) *log.Logger {
	l := log.New(os.Stdout, "["+scope+"] ", log.Lmsgprefix|log.LstdFlags)
	return l
}

var IpcLogger = NewLogger("ipc")
var AppLogger = NewLogger("apero")

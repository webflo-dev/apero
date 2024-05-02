package utils

import (
	"log"
	"os"
)

func NewLogger(scope string) *log.Logger {
	l := log.New(os.Stdout, "["+scope+"] ", log.Lmsgprefix|log.LstdFlags)
	return l
}

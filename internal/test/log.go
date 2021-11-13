package test

import (
	"io"
	"log"
)

func NewLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

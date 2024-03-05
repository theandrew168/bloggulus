package task

import (
	"log"
	"sync"
)

type Worker struct {
	sync.WaitGroup
	logger *log.Logger
}

func NewWorker(logger *log.Logger) *Worker {
	worker := Worker{
		logger: logger,
	}
	return &worker
}

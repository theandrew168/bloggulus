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

func (w *Worker) log(msg string) {
	w.logger.Output(2, msg)
}

func (w *Worker) logError(err error) {
	w.logger.Output(2, err.Error())
}

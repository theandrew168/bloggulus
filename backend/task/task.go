package task

import (
	"time"
)

type Task interface {
	Run(interval time.Duration)
	RunNow() error
}

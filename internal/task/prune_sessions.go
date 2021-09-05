package task

import (
	"context"
	"log"
	"time"

	"github.com/theandrew168/bloggulus/internal/core"
)

type pruneSessionsTask struct {
	Session core.SessionStorage
}

func PruneSessions(sessionStorage core.SessionStorage) Task {
	return &pruneSessionsTask{
		Session: sessionStorage,
	}
}

func (t *pruneSessionsTask) Run(interval time.Duration) {
	c := time.Tick(interval)
	for {
		<-c

		err := t.pruneSessions()
		if err != nil {
			log.Println(err)
		}
	}
}

func (t *pruneSessionsTask) RunNow() error {
	return t.pruneSessions()
}

func (t *pruneSessionsTask) pruneSessions() error {
	return t.Session.DeleteExpired(context.Background())
}

package task

import (
	"context"
	"log"
	"time"

	"github.com/theandrew168/bloggulus/model"
)

type pruneSessionsTask struct {
	Session model.SessionStorage
}

func PruneSessions(sessionStorage model.SessionStorage) Task {
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

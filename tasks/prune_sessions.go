package tasks

import (
	"context"
	"log"
	"time"

	"github.com/theandrew168/bloggulus/storage"
)

type pruneSessionsTask struct {
	Session storage.Session
}

func PruneSessions(sessionStorage storage.Session) Task {
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

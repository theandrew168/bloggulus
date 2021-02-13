package tasks

import (
	"context"
	"log"
	"time"

	"github.com/theandrew168/bloggulus/storage"
)

type cleanupSessionsTask struct {
	Session storage.Session
}

func CleanupSessions(sessionStorage storage.Session) Task {
	return &cleanupSessionsTask{
		Session: sessionStorage,
	}
}

func (t *cleanupSessionsTask) Run(interval time.Duration) {
	c := time.Tick(interval)
	for {
		<-c

		err := t.cleanupSessions()
		if err != nil {
			log.Println(err)
		}
	}
}

func (t *cleanupSessionsTask) RunNow() error {
	return t.cleanupSessions()
}

func (t *cleanupSessionsTask) cleanupSessions() error {
	return t.Session.DeleteExpired(context.Background())
}

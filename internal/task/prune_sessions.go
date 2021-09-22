package task

import (
	"context"
	"log"
	"time"

	"github.com/theandrew168/bloggulus/internal/core"
)

type pruneSessionsTask struct {
	session core.SessionStorage
}

func PruneSessions(session core.SessionStorage) Task {
	return &pruneSessionsTask{
		session: session,
	}
}

func (t *pruneSessionsTask) Run(interval time.Duration) {
	err := t.RunNow()
	if err != nil {
		log.Println(err)
	}

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
	return t.session.DeleteExpired(context.Background())
}

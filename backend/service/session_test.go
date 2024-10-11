package service_test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestClearExpiredSessions(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	account := test.CreateAccount(t, repo)

	sessionOld, _, err := model.NewSession(
		account,
		-1*time.Hour,
	)
	test.AssertNilError(t, err)

	err = repo.Session().Create(sessionOld)
	test.AssertNilError(t, err)

	sessionNew, _, err := model.NewSession(
		account,
		1*time.Hour,
	)
	test.AssertNilError(t, err)

	err = repo.Session().Create(sessionNew)
	test.AssertNilError(t, err)

	s := service.NewSessionService(repo)
	err = s.ClearExpiredSessions()
	test.AssertNilError(t, err)

	_, err = repo.Session().Read(sessionOld.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)

	_, err = repo.Session().Read(sessionNew.ID())
	test.AssertNilError(t, err)
}

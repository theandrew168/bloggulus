package command_test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestDeleteExpiredSessions(t *testing.T) {
	t.Parallel()

	repo, closer := test.NewRepository(t)
	defer closer()

	cmd := command.NewAuth(repo)

	account := test.CreateAccount(t, repo)

	sessionOld, _, err := model.NewSession(model.NewSessionParams{
		Account: account,
		TTL:     -1 * time.Hour,
	})
	test.AssertNilError(t, err)

	err = repo.Session().Create(sessionOld)
	test.AssertNilError(t, err)

	sessionNew, _, err := model.NewSession(model.NewSessionParams{
		Account: account,
		TTL:     1 * time.Hour,
	})
	test.AssertNilError(t, err)

	err = repo.Session().Create(sessionNew)
	test.AssertNilError(t, err)

	err = cmd.DeleteExpiredSessions(timeutil.Now())
	test.AssertNilError(t, err)

	_, err = repo.Session().Read(sessionOld.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)

	_, err = repo.Session().Read(sessionNew.ID())
	test.AssertNilError(t, err)
}

package config_test

import (
	"fmt"
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/testutil"
)

const (
	databaseURI = "postgresql://foo:bar@localhost:5432/postgres"
	port        = "5000"
)

func TestRead(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		database_uri = "%s"
		port = "%s"
	`, databaseURI, port)

	cfg, err := config.Read(data)
	testutil.AssertNilError(t, err)

	testutil.AssertEqual(t, cfg.DatabaseURI, databaseURI)
	testutil.AssertEqual(t, cfg.Port, port)
}

func TestOptional(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		database_uri = "%s"
	`, databaseURI)

	cfg, err := config.Read(data)
	testutil.AssertNilError(t, err)

	testutil.AssertEqual(t, cfg.DatabaseURI, databaseURI)
	testutil.AssertEqual(t, cfg.Port, config.DefaultPort)
}

func TestRequired(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		port = "%s"
	`, port)

	_, err := config.Read(data)
	testutil.AssertErrorContains(t, err, "missing")
	testutil.AssertErrorContains(t, err, "database_uri")
}

func TestExtra(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		database_uri = "%s"
		foo = "bar"
	`, databaseURI)

	_, err := config.Read(data)
	testutil.AssertErrorContains(t, err, "extra")
	testutil.AssertErrorContains(t, err, "foo")
}

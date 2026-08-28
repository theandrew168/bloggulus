package config_test

import (
	"fmt"
	"testing"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/test"
)

const (
	databaseURI = "postgresql://foo:bar@localhost:5432/postgres"
	secretKey   = "secret"
)

func TestRead(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		database_uri = "%s"
		secret_key = "%s"
	`, databaseURI, secretKey)

	cfg, err := config.Read(data)
	test.AssertNilError(t, err)

	test.AssertEqual(t, cfg.DatabaseURI, databaseURI)
}

func TestOptional(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		database_uri = "%s"
		secret_key = "%s"
	`, databaseURI, secretKey)

	cfg, err := config.Read(data)
	test.AssertNilError(t, err)

	test.AssertEqual(t, cfg.DatabaseURI, databaseURI)
}

func TestRequired(t *testing.T) {
	t.Parallel()

	data := ""

	_, err := config.Read(data)
	test.AssertErrorContains(t, err, "missing")
	test.AssertErrorContains(t, err, "database_uri")
	test.AssertErrorContains(t, err, "secret_key")
}

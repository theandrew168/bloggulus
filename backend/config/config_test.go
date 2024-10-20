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
	port        = "5000"
)

func TestRead(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		database_uri = "%s"
		secret_key = "%s"
		port = "%s"
	`, databaseURI, secretKey, port)

	cfg, err := config.Read(data)
	test.AssertNilError(t, err)

	test.AssertEqual(t, cfg.DatabaseURI, databaseURI)
	test.AssertEqual(t, cfg.Port, port)
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
	test.AssertEqual(t, cfg.Port, config.DefaultPort)
}

func TestRequired(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		port = "%s"
	`, port)

	_, err := config.Read(data)
	test.AssertErrorContains(t, err, "missing")
	test.AssertErrorContains(t, err, "database_uri")
	test.AssertErrorContains(t, err, "secret_key")
}

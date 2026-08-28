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

func TestOAuth(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		database_uri = "%s"
		secret_key = "%s"
		github_client_id = "github_client_id"
		github_client_secret = "github_client_secret"
		github_redirect_uri = "github_redirect_uri"
		google_client_id = "google_client_id"
		google_client_secret = "google_client_secret"
		google_redirect_uri = "google_redirect_uri"
	`, databaseURI, secretKey)

	cfg, err := config.Read(data)
	test.AssertNilError(t, err)

	test.AssertEqual(t, cfg.GithubOAuth.ClientID, "github_client_id")
	test.AssertEqual(t, cfg.GithubOAuth.ClientSecret, "github_client_secret")
	test.AssertEqual(t, cfg.GithubOAuth.RedirectURI, "github_redirect_uri")
	test.AssertEqual(t, cfg.GoogleOAuth.ClientID, "google_client_id")
	test.AssertEqual(t, cfg.GoogleOAuth.ClientSecret, "google_client_secret")
	test.AssertEqual(t, cfg.GoogleOAuth.RedirectURI, "google_redirect_uri")
}

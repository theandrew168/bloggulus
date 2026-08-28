package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type OAuthConfig struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURI  string `toml:"redirect_uri"`
}

type Config struct {
	DatabaseURI     string
	SecretKey       string
	EnableDebugAuth bool
	GithubOAuth     OAuthConfig
	GoogleOAuth     OAuthConfig
}

type configFile struct {
	DatabaseURI        string `toml:"database_uri"`
	SecretKey          string `toml:"secret_key"`
	EnableDebugAuth    bool   `toml:"enable_debug_auth"`
	GithubClientID     string `toml:"github_client_id"`
	GithubClientSecret string `toml:"github_client_secret"`
	GithubRedirectURI  string `toml:"github_redirect_uri"`
	GoogleClientID     string `toml:"google_client_id"`
	GoogleClientSecret string `toml:"google_client_secret"`
	GoogleRedirectURI  string `toml:"google_redirect_uri"`
}

func Read(data string) (Config, error) {
	// Initialize config with default values (if applicable).
	file := configFile{
		EnableDebugAuth: false,
	}
	meta, err := toml.Decode(data, &file)
	if err != nil {
		return Config{}, err
	}

	// Build set of present config keys.
	present := make(map[string]bool)
	for _, keys := range meta.Keys() {
		key := keys[0]
		present[key] = true
	}

	required := []string{
		"database_uri",
		"secret_key",
	}

	// Gather any missing values.
	missing := []string{}
	for _, key := range required {
		if _, ok := present[key]; !ok {
			missing = append(missing, key)
		}
	}

	// Error upon missing values
	if len(missing) > 0 {
		msg := strings.Join(missing, ", ")
		return Config{}, fmt.Errorf("missing config values: %s", msg)
	}

	conf := Config{
		DatabaseURI:     file.DatabaseURI,
		SecretKey:       file.SecretKey,
		EnableDebugAuth: file.EnableDebugAuth,
	}

	if file.GithubClientID != "" && file.GithubClientSecret != "" && file.GithubRedirectURI != "" {
		conf.GithubOAuth = OAuthConfig{
			ClientID:     file.GithubClientID,
			ClientSecret: file.GithubClientSecret,
			RedirectURI:  file.GithubRedirectURI,
		}
	}

	if file.GoogleClientID != "" && file.GoogleClientSecret != "" && file.GoogleRedirectURI != "" {
		conf.GoogleOAuth = OAuthConfig{
			ClientID:     file.GoogleClientID,
			ClientSecret: file.GoogleClientSecret,
			RedirectURI:  file.GoogleRedirectURI,
		}
	}

	return conf, nil
}

func ReadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	return Read(string(data))
}

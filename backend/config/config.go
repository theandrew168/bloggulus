package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	DefaultPort = "5000"
)

type Config struct {
	DatabaseURI        string `toml:"database_uri"`
	Port               string `toml:"port"`
	GithubClientID     string `toml:"github_client_id"`
	GithubClientSecret string `toml:"github_client_secret"`
	GoogleClientID     string `toml:"google_client_id"`
	GoogleClientSecret string `toml:"google_client_secret"`
}

func Read(data string) (Config, error) {
	// Initialize config with default values.
	conf := Config{
		Port: DefaultPort,
	}
	meta, err := toml.Decode(data, &conf)
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

	return conf, nil
}

func ReadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	return Read(string(data))
}

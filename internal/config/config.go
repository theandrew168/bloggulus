package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	defaultPort = "5000"
)

type Config struct {
	DatabaseURI string `toml:"database_uri"`
	Port        string `toml:"port"`
}

func Read(data string) (Config, error) {
	var cfg Config
	meta, err := toml.Decode(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	// gather extra values
	extra := []string{}
	for _, keys := range meta.Undecoded() {
		key := keys[0]
		extra = append(extra, key)
	}

	// error upon extra values
	if len(extra) > 0 {
		msg := strings.Join(extra, ", ")
		return Config{}, fmt.Errorf("extra config values: %s", msg)
	}

	// build set of present config keys
	present := make(map[string]bool)
	for _, keys := range meta.Keys() {
		key := keys[0]
		present[key] = true
	}

	required := []string{
		"database_uri",
	}

	// gather missing values
	missing := []string{}
	for _, key := range required {
		if _, ok := present[key]; !ok {
			missing = append(missing, key)
		}
	}

	// error upon missing values
	if len(missing) > 0 {
		msg := strings.Join(missing, ", ")
		return Config{}, fmt.Errorf("missing config values: %s", msg)
	}

	// handle defaults
	if cfg.Port == "" {
		cfg.Port = defaultPort
	}

	return cfg, nil
}

func ReadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	return Read(string(data))
}

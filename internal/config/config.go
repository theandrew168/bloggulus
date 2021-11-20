package config

import (
	_ "embed"

	"github.com/BurntSushi/toml"
)

//go:embed defaults.toml
var defaults string

type Config struct {
	Env         string `toml:"env"`
	Port        string `toml:"port"`
	DatabaseURL string `toml:"database_url"`
}

func Defaults() Config {
	var cfg Config
	_, err := toml.Decode(defaults, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func FromFile(path string) (Config, error) {
	// start with defaults
	defaults := Defaults()

	// attempt to read user-defined config file
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, err
	}

	// merge the results
	// TODO: is there a way to do this automatically? with reflect?
	if cfg.Env == "" {
		cfg.Env = defaults.Env
	}
	if cfg.Port == "" {
		cfg.Port = defaults.Port
	}
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = defaults.DatabaseURL
	}

	return cfg, nil
}

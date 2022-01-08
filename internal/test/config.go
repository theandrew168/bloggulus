package test

import (
	_ "embed"
	"testing"

	"github.com/theandrew168/bloggulus/internal/config"
)

func Config(t *testing.T) config.Config {
	t.Helper()

	// read the local development config file
	cfg, err := config.ReadFile("../../bloggulus.conf")
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

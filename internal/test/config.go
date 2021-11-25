package test

import (
	_ "embed"
	"testing"

	"github.com/theandrew168/bloggulus/internal/config"
)

//go:embed bloggulus.conf
var testConfig string

func Config(t *testing.T) config.Config {
	t.Helper()

	cfg, err := config.Read(testConfig)
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

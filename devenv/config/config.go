package config

import (
	"github.com/felipeweb/devctl/devenv/context"
)

type Config struct {
	*context.Context
}

// New Config initializes a default porter configuration.
func New() *Config {
	return &Config{
		Context: context.New(),
	}
}

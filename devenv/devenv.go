package devenv

import (
	"github.com/felipeweb/devctl/devenv/config"
)

// Devenv is the logic behind the devctl client.
type Devenv struct {
	*config.Config
}

// New devenv client, initialized with useful defaults.
func New() *Devenv {
	return &Devenv{
		Config: config.New(),
	}
}

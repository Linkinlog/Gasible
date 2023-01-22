package gitService

import (
	"errors"

	"github.com/Linkinlog/gasible/internal/models"
)

type Opts struct {
	NoOp       bool
	CreateRepo bool
}

func (o Opts) Run(g *models.GitServiceConfig) ([]byte, error) {
	// Set up Git options
	// Set up / pull bare repo
	// Sync repo and sensibly overwrite defaults
	return []byte{}, errors.New("Needs implemented!")
}

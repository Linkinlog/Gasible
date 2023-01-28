package gitService

import (
	"errors"

	"github.com/Linkinlog/gasible/internal/models"
)

type Opts struct {
	NoOp bool
}

func (o Opts) Run(g *models.GitServiceConfig) ([][]byte, error) {
	if g.RepoName == "" {
		return [][]byte{}, errors.New("No RepoName defined. Needed for bare repo setup.")
	}
	output := [][]byte{}
	os := models.GetCurrentSystem()
	// Set up Git options
	out, err := g.Setup(o.NoOp, *os)
	if err := handleLog(&output, out, err); err != nil {
		return [][]byte{}, err
	}
	// Set up the alias for configuring the repo
	out, err = g.SetupRepoAlias(o.NoOp, *os)
	if err := handleLog(&output, out, err); err != nil {
		return [][]byte{}, err
	}
	// Sync repo and sensibly overwrite defaults
	return output, nil
}

func handleLog(output *[][]byte, out []byte, err error) error {
	if err != nil {
		return err
	}
	*output = append(*output, out)
	return nil
}

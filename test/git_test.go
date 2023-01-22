package test

import (
	"testing"

	"github.com/Linkinlog/gasible/cmd/gitService"
	"github.com/Linkinlog/gasible/internal/models"
)

func TestGitSetup(t *testing.T) {
	gitConfig := models.GitServiceConfig{
		Enabled: true,
		Options: &models.Options{
			User:     "Linkinlog",
			Email:    "",
			RepoURL:  "https://github.com/Linkinlog/.dotfiles.git",
			RepoName: "",
		},
	}
	gitServiceOpts := gitService.Opts{
		NoOp:       false,
		CreateRepo: false,
	}
	_, err := gitServiceOpts.Run(&gitConfig)
	if err != nil {
		t.Fatal(err)
	}
}

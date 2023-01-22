package test

import (
	"testing"

	"github.com/Linkinlog/gasible/cmd/gitService"
	"github.com/Linkinlog/gasible/internal/models"
)

func TestGitSetup(t *testing.T) {
	gitOpts := models.GitServiceConfig{
		Enabled: true,
		Options: &models.Options{
			User:     "Linkinlog",
			Email:    "dahlton@dahlton.org",
			RepoURL:  "https://github.com/Linkinlog/.dotfiles.git",
			RepoName: "",
		},
	}
	gitServiceOpts := gitService.Opts{
		NoOp:       false,
		CreateRepo: false,
	}
	_, err := gitServiceOpts.Run(&gitOpts)
	if err != nil {
		t.Fatal(err)
	}
}

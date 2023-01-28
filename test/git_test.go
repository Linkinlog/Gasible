package test

import (
	"testing"

	"github.com/Linkinlog/gasible/cmd/gitService"
	"github.com/Linkinlog/gasible/internal/models"
)

var (
	gitConfig = models.GitServiceConfig{
		Enabled: true,
		Options: models.Options{
			User:     "Linkinlog",
			Email:    "",
			RepoURL:  "https://github.com/Linkinlog/.dotfiles.git",
			RepoName: "",
		},
	}
	gitServiceOpts = gitService.Opts{
		NoOp:       false,
		CreateRepo: false,
	}
)

func TestGitRun(t *testing.T) {
	_, err := gitServiceOpts.Run(&gitConfig)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGitSetup(t *testing.T) {
	_, err := gitConfig.Setup(true)
	if err != nil {
		t.Fatal(err)
	}
}

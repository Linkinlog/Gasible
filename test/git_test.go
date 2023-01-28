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
			Create:   false,
			RepoURL:  "https://github.com/Linkinlog/.dotfiles.git",
			RepoName: "",
		},
	}
	gitServiceOpts = gitService.Opts{
		NoOp: false,
	}
)

func TestGitRun(t *testing.T) {
	out, err := gitServiceOpts.Run(&gitConfig)
	if err != nil || len(out) < 1 {
		t.Fatal(err)
	}
}

func TestGitSetup(t *testing.T) {
	os := models.GetCurrentSystem()
	out, err := gitConfig.Setup(true, *os)
	if err != nil || len(out) < 1 {
		t.Fatal(err)
	}
}

func TestGitSetupRepo(t *testing.T) {
	out, err := gitConfig.SetupRepo(true)
	if err != nil || len(out) < 1 {
		t.Fatal(err)
	}
}

func TestGitSetupRepoAlias(t *testing.T) {
	os := models.GetCurrentSystem()
	out, err := gitConfig.SetupRepoAlias(true, *os)
	if err != nil || len(out) < 1 {
		t.Fatal(err)
	}
}

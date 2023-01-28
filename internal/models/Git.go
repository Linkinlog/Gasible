package models

import (
	"fmt"
	"os/exec"
)

type GitServiceConfig struct {
	Enabled bool `yaml:"Enabled"`
	Options `yaml:"Options"`
}
type Options struct {
	User     string `yaml:"User"`
	Email    string `yaml:"Email"`
	Create   bool   `yaml:"Create Repo"`
	RepoURL  string `yaml:"Repo Clone URL"`
	RepoName string `yaml:"Repo Name"`
}

func (g GitServiceConfig) Default() *GitServiceConfig {
	return &GitServiceConfig{
		Enabled: true,
		Options: Options{},
	}
}

// Setup sets up some Git config options.
func (g GitServiceConfig) Setup(NoOp bool, os System) ([]byte, error) {
	// git config --global user.name "John Doe"
	userNameConfig := g.Options.User
	command := "git config --global user.name \"" + userNameConfig + "\""
	out, err := os.Exec(NoOp, command)
	if err != nil {
		return []byte{}, err
	}
	// git config --global user.email johndoe@example.com
	emailConfig := g.Options.Email
	command = "git config --global user.email \"" + emailConfig + "\""
	emailOut, err := os.Exec(NoOp, command)
	if err != nil {
		return []byte{}, err
	}
	out = []byte(string(out) + "\n" + string(emailOut))
	return out, nil
}

// Create a bare git repo using the RepoName.
func (g GitServiceConfig) SetupRepo(NoOp bool) ([]byte, error) {
	// TODO
	return []byte{}, nil
}

// SetupRepoAlias creates an alias for us to interact with our bare repo.
// Requires the RepoName to be set in order to be effective.
func (g GitServiceConfig) SetupRepoAlias(NoOp bool, os System) ([]byte, error) {
	// Find the full path to wherever git is
	gitExec, err := exec.LookPath("git")
	if err != nil {
		return []byte{}, err
	}
	// Format & Run the alias command
	command := fmt.Sprintf("alias config=\"%s --git-dir=$HOME/%s --work-tree=$HOME\"", gitExec, g.RepoName)
	out, err := os.Exec(NoOp, command)
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

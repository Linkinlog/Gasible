package models

type GitServiceConfig struct {
	Enabled  bool `yaml:"Enabled"`
	*Options `yaml:"Options"`
}
type Options struct {
	User     string `yaml:"User"`
	Email    string `yaml:"Email"`
	RepoURL  string `yaml:"Repo Clone URL"`
	RepoName string `yaml:"Repo Name"`
}

func (g GitServiceConfig) Default() *GitServiceConfig {
	return &GitServiceConfig{
		Enabled: false,
		Options: &Options{},
	}
}

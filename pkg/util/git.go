package util

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
)

// GetRepositoryRemotes returns a map containing the remote names and the URLs they
// point to.
func GetRepositoryRemotes(repo *git.Repository) (map[string][]string, error) {
	gitRemotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("cannot list remotes: %v", err)
	}

	remotes := map[string][]string{}
	for _, remote := range gitRemotes {
		remoteCfg := remote.Config()
		remotes[remoteCfg.Name] = remoteCfg.URLs
	}
	return remotes, nil
}

// UpdateRepositoryRemotes updates the URLs list for a specific remote.
func UpdateRepositoryRemotes(repo *git.Repository, name string, URLs []string) error {
	cfg, err := repo.Storer.Config()
	if err != nil {
		return err
	}

	cfg.Remotes[name] = &gitconfig.RemoteConfig{
		Name: name,
		URLs: URLs,
	}
	return repo.Storer.SetConfig(cfg)
}

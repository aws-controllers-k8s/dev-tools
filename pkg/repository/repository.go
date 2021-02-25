package repository

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
)

// NewRepository returns a pointer to a new repository.
func NewRepository(name string, repoType RepositoryType) *Repository {
	if repoType == RepositoryTypeController {
		name = fmt.Sprintf("%s-controller", name)
	}
	return &Repository{
		Name: name,
		Type: repoType,
	}
}

// Repository represents an ACK project repository.
type Repository struct {
	gitRepo *git.Repository

	Name        string
	FullPath    string
	GitHead     string
	ForkName    string
	ForkURL     string
	UpstreamURL string
	Type        RepositoryType
	RemoteURL   string
}

func remoteURL(owner, name string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", owner, name)
}

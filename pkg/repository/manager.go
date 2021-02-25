package repository

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	git "gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/github"
	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
)

const (
	originRemoteName   = "origin"
	upstreamRemoteName = "upstream"

	defaultRemoteName = originRemoteName
)

var (
	ErrorSSHKeyNotFound           error = errors.New("ssh key not found")
	ErrorRepositoryUnknown        error = errors.New("unknown repository")
	ErrorRepositoryAlreadyExist   error = errors.New("repository already exist")
	ErrorMissingGithubCredentials error = errors.New("missing github credentials")
)

// NewManager create a new manager
func NewManager(cfg *config.Config) (*Manager, error) {
	ghc := github.NewClient(cfg.Github.Token)

	var sshKey []byte
	if cfg.Git.SSHKeyPath != "" {
		sshKeyBytes, err := ioutil.ReadFile(cfg.Git.SSHKeyPath)
		if err != nil {
			return nil, ErrorSSHKeyNotFound
		}
		sshKey = sshKeyBytes
	}

	return &Manager{
		log:    logrus.New(),
		cfg:    cfg,
		sshKey: sshKey,
		ghc:    ghc,
	}, nil
}

// Manager is reponsible of managing local ACK local repositories and
// github forks.
type Manager struct {
	log       *logrus.Logger
	cfg       *config.Config
	ghc       *github.Client
	sshKey    []byte
	repoCache []*Repository
}

// LoadAll parses the configuration and loads informations about local
// repositories if they are found.
func (m *Manager) LoadAll() error {
	// collect repositories from config
	for _, coreRepo := range m.cfg.Repositories.Core {
		repo, err := m.LoadRepository(coreRepo, RepositoryTypeCore)
		if err != nil {
			return err
		}
		m.repoCache = append(m.repoCache, repo)
	}
	for _, serviceName := range m.cfg.Repositories.Services {
		repo, err := m.LoadRepository(serviceName, RepositoryTypeController)
		if err != nil {
			return err
		}
		m.repoCache = append(m.repoCache, repo)
	}
	return nil
}

// LoadRepository loads informations a single repository
func (m *Manager) LoadRepository(name string, t RepositoryType) (*Repository, error) {
	// check repo cache
	repo, err := m.GetRepository(name)
	if err == nil {
		return repo, nil
	}

	switch t {
	case RepositoryTypeCore:
		if !util.InStrings(name, m.cfg.Repositories.Core) {
			return nil, fmt.Errorf("core not found")
		}
	case RepositoryTypeController:
		if !util.InStrings(name, m.cfg.Repositories.Services) {
			return nil, fmt.Errorf("controller not configured")
		}
	}

	// controller repositories should always have a '-controller' suffix
	if t == RepositoryTypeController {
		name = fmt.Sprintf("%s-controller", name)
	}

	// set fork name
	forkName := name
	if m.cfg.Github.ForkPrefix != "" {
		forkName = fmt.Sprintf("%s%s", m.cfg.Github.ForkPrefix, name)
	}

	fullPath := filepath.Join(m.cfg.RootDirectory, name)
	var gitRepo *git.Repository
	var gitHead string
	_, err = os.Stat(fullPath)
	if !os.IsNotExist(err) {
		gitRepo, err = git.PlainOpen(fullPath)
		if err != nil {
			return nil, err
		}
		head, err := gitRepo.Head()
		if err != nil {
			return nil, err
		}
		gitHead = head.Name().Short()
	}

	return &Repository{
		Name:      name,
		Type:      t,
		gitRepo:   gitRepo,
		GitHead:   gitHead,
		FullPath:  fullPath,
		RemoteURL: remoteURL(m.cfg.Github.Username, forkName),
	}, nil
}

// HasRepo return true if repoName is known
func (m *Manager) hasRepo(repoName string) bool {
	for _, repo := range m.repoCache {
		if repo.Name == repoName {
			return true
		}
	}
	return false
}

// GetRepository return a known repository
func (m *Manager) GetRepository(repoName string) (*Repository, error) {
	for _, repo := range m.repoCache {
		if repo.Name == repoName {
			return repo, nil
		}
	}
	return nil, ErrorRepositoryUnknown
}

// List returns a list of filtered repositories
func (m *Manager) List(filters ...Filter) []*Repository {
	return m.ListAnd(filters...)
}

// List returns a list of filtered repositories
func (m *Manager) ListAnd(filters ...Filter) []*Repository {
	repos := []*Repository{}
mainLoop:
	for _, repo := range m.repoCache {
		for _, filter := range filters {
			if !filter(repo) {
				continue mainLoop
			}
		}
		repos = append(repos, repo)
	}
	return repos
}

// List returns a list of filtered repositories
func (m *Manager) ListOr(filters ...Filter) []*Repository {
	repos := []*Repository{}
mainLoop:
	for _, repo := range m.repoCache {
		for _, filter := range filters {
			if filter(repo) {
				repos = append(repos, repo)
				continue mainLoop
			}
		}
	}
	return repos
}

// Clone clones a known repository to the root directory
func (m *Manager) clone(repoName string) error {
	repo, err := m.GetRepository(repoName)
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %v", repoName, err)
	}
	if repo.gitRepo != nil {
		return ErrorRepositoryAlreadyExist
	}

	signer, err := ssh.ParsePrivateKey([]byte(m.sshKey))
	if err != nil {
		return err
	}

	// Clone repository
	auth := &gitssh.PublicKeys{User: "git", Signer: signer}
	gitRepo, err := git.PlainClone(repo.FullPath, false, &git.CloneOptions{
		Auth:       auth,
		URL:        repo.ForkURL,
		RemoteName: defaultRemoteName,
		Progress:   nil,
	})
	if errors.Is(err, transport.ErrAuthenticationRequired) {
		return ErrorMissingGithubCredentials
	}
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %v", repoName, err)
	}
	repo.gitRepo = gitRepo

	// Add upstream remote
	gitRepo.CreateRemote(&gitconfig.RemoteConfig{
		Name: upstreamRemoteName,
		URLs: []string{},
	})

	return nil
}

// ensureFork ensures that your github account have a fork for a given
// project. It will also rename the project if it's not following the
// standard: $ackprefix-$projectname
func (m *Manager) ensureFork(repo *Repository) error {
	ctx := context.TODO()
	m.log.SetLevel(logrus.DebugLevel)

	expectedForkName := fmt.Sprintf("%s%s", m.cfg.Github.ForkPrefix, repo.Name)
	fork, err := m.ghc.GetRepositoryForkForUser(ctx, repo.Name, m.cfg.Github.Username)
	if err == nil {
		if fork.Name != expectedForkName {
			err = m.ghc.RenameRepository(ctx, m.cfg.Github.Username, fork.Name, expectedForkName)
			if err != nil {
				return err
			}
			repo.ForkName = expectedForkName
		}
	} else if err == github.ErrorForkNotFound {
		err = m.ghc.ForkRepository(ctx, repo.Name)
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Second)

		err = m.ghc.RenameRepository(ctx, m.cfg.Github.Username, repo.Name, expectedForkName)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (m *Manager) ensureClone(repo *Repository) error {
	err := m.clone(repo.Name)
	if err != nil && err != ErrorRepositoryAlreadyExist {
		return err
	}
	return nil
}

// EnsureAll ensures all repositories
func (m *Manager) EnsureAll() error {
	for _, repo := range m.repoCache {
		err := m.ensureFork(repo)
		if err != nil {
			return err
		}

		err = m.ensureClone(repo)
		if err != nil {
			return err
		}
	}
	return nil
}

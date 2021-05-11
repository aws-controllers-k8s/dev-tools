// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	ackdevgit "github.com/aws-controllers-k8s/dev-tools/pkg/git"
	"github.com/aws-controllers-k8s/dev-tools/pkg/github"
	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
)

const (
	originRemoteName   = "origin"
	upstreamRemoteName = "upstream"
)

var (
	ErrUnauthenticated        error = errors.New("unauthenticated")
	ErrUnconfiguredRepository error = errors.New("unconfigured repository")
	ErrRepositoryNotCached    error = errors.New("repository not cached")
	ErrRepositoryDoesntExist  error = errors.New("repository doesnt exist")
	ErrRepositoryAlreadyExist error = errors.New("repository already exist")
)

// NewManager create a new manager.
func NewManager(cfg *config.Config) (*Manager, error) {
	githubClient := github.NewClient(cfg.Github.Token)
	gitOpts := []ackdevgit.Option{
		ackdevgit.WithRemote(originRemoteName),
	}
	urlBuilder := httpsRemoteURL

	// Add git authentication options
	if cfg.Git.SSHKeyPath == "" {
		gitOpts = append(gitOpts,
			ackdevgit.WithGithubCredentials(cfg.Github.Username, cfg.Github.Token),
		)
	} else {
		// TODO(hilalymh) set ssh.Signer here.. figure out how to deal with encrypted
		// keys properly...
		gitOpts = append(gitOpts, ackdevgit.WithSSHSigner(nil))
		urlBuilder = sshRemoteURL
	}

	gitClient := ackdevgit.New(gitOpts...)

	return &Manager{
		repoCache: make(map[string]*Repository),

		cfg:        cfg,
		ghc:        githubClient,
		git:        gitClient,
		urlBuilder: urlBuilder,
	}, nil
}

// Manager is reponsible of managing local ACK local repositories and
// github forks.
type Manager struct {
	repoCache map[string]*Repository

	log        *logrus.Logger
	cfg        *config.Config
	git        ackdevgit.OpenCloner
	ghc        github.RepositoryService
	urlBuilder func(owner, repo string) string
}

// LoadRepository loads information about a single local repository
func (m *Manager) LoadRepository(name string, t RepositoryType) (*Repository, error) {
	// check repo cache
	repo, err := m.GetRepository(name)
	if err == nil {
		return repo, nil
	}

	// fail if repository doesn't exist in the manager configuration
	switch t {
	case RepositoryTypeCore:
		if !util.InStrings(name, m.cfg.Repositories.Core) {
			return nil, ErrUnconfiguredRepository
		}
	case RepositoryTypeController:
		if !util.InStrings(name, m.cfg.Repositories.Services) {
			return nil, ErrUnconfiguredRepository
		}
	}

	repoName := name
	// controller repositories should always have a '-controller' suffix
	if t == RepositoryTypeController {
		repoName = fmt.Sprintf("%s-controller", name)
	}

	// set expected fork name
	expectedForkName := repoName
	if m.cfg.Github.ForkPrefix != "" {
		expectedForkName = fmt.Sprintf("%s%s", m.cfg.Github.ForkPrefix, repoName)
	}

	var gitHead string
	var gitRepo *git.Repository
	fullPath := filepath.Join(m.cfg.RootDirectory, repoName)

	gitRepo, err = m.git.Open(fullPath)
	if err != nil && err != git.ErrRepositoryNotExists {
		return nil, err
	} else if err == nil {
		// load current branch
		head, err := gitRepo.Head()
		if err != nil {
			return nil, err
		}
		gitHead = head.Name().Short()
	}

	repo = &Repository{
		Name:             repoName,
		Type:             t,
		gitRepo:          gitRepo,
		GitHead:          gitHead,
		FullPath:         fullPath,
		ExpectedForkName: expectedForkName,
	}

	// cache repository
	m.repoCache[name] = repo
	return repo, nil
}

// LoadAll parses the configuration and loads informations about local
// repositories if they are found.
func (m *Manager) LoadAll() error {
	// collect repositories from config
	for _, coreRepo := range m.cfg.Repositories.Core {
		_, err := m.LoadRepository(coreRepo, RepositoryTypeCore)
		if err != nil {
			return err
		}
	}
	for _, serviceName := range m.cfg.Repositories.Services {
		_, err := m.LoadRepository(serviceName, RepositoryTypeController)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetRepository return a known repository
func (m *Manager) GetRepository(repoName string) (*Repository, error) {
	repo, ok := m.repoCache[repoName]
	if !ok {
		return nil, ErrRepositoryNotCached
	}
	return repo, nil
}

// List returns the list of all the cached repositories
func (m *Manager) List(filters ...Filter) []*Repository {
	repos := []*Repository{}
	repoNames := append(m.cfg.Repositories.Core, m.cfg.Repositories.Services...)
mainLoop:
	for _, repoName := range repoNames {
		repo, err := m.GetRepository(repoName)
		if err != nil {
			continue
		}
		for _, filter := range filters {
			if !filter(repo) {
				continue mainLoop
			}
		}
		repos = append(repos, repo)
	}
	return repos
}

// Clone clones a known repository to the config root directory
func (m *Manager) clone(ctx context.Context, repoName string) error {
	repo, err := m.GetRepository(repoName)
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %v", repoName, err)
	}
	if repo.gitRepo != nil {
		return ErrRepositoryAlreadyExist
	}

	// clone fork repository
	err = m.git.Clone(
		ctx,
		m.urlBuilder(m.cfg.Github.Username, repo.ExpectedForkName),
		repo.ExpectedForkName,
	)
	if errors.Is(err, transport.ErrAuthenticationRequired) {
		return ErrUnauthenticated
	}
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %v", repoName, err)
	}

	// open git repository
	gitRepo, err := m.git.Open(repo.FullPath)
	if err != nil {
		// Maybe panic here?
		return err
	}

	// set repository git object
	repo.gitRepo = gitRepo

	// Add upstream remote
	_, err = gitRepo.CreateRemote(&gitconfig.RemoteConfig{
		Name: upstreamRemoteName,
		URLs: []string{m.urlBuilder(github.ACKOrg, repo.Name)},
	})

	if err != nil {
		return fmt.Errorf("cannot add upstream remote to repository %s: %v", repoName, err)
	}

	return nil
}

// ensureFork ensures that your github account have a fork for a given
// ACK project. It will also rename the project if it's not following the
// standard: $ackprefix-$projectname
func (m *Manager) EnsureFork(ctx context.Context, repo *Repository) error {
	// TODO(hilaly): m.log.SetLevel(logrus.DebugLevel)

	fork, err := m.ghc.GetUserRepositoryFork(ctx, repo.Name)
	if err == nil {
		if *fork.Name != repo.ExpectedForkName {
			err = m.ghc.RenameRepository(ctx, m.cfg.Github.Username, *fork.Name, repo.ExpectedForkName)
			if err != nil {
				return err
			}
		}
	} else if err == github.ErrorForkNotFound {
		err = m.ghc.ForkRepository(ctx, repo.Name)
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Second)

		err = m.ghc.RenameRepository(ctx, m.cfg.Github.Username, repo.Name, repo.ExpectedForkName)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

// used to help mocking os.Rename
// TODO(hilalymh): Q4 switch to go1.16 os/fs library/interface
var renameDirectory = os.Rename

func (m *Manager) EnsureClone(ctx context.Context, repo *Repository) error {
	err := m.clone(ctx, repo.Name)
	if err != nil && err != ErrRepositoryAlreadyExist {
		return err
	}

	// At this point we ensured that the fork repository is cloned. We need to rename it
	// if there is any fork prefix.
	if repo.Name != repo.ExpectedForkName {

		newPath := filepath.Join(
			filepath.Dir(repo.FullPath),
			repo.Name,
		)
		err := renameDirectory(repo.FullPath, newPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// EnsureAll ensures one repository.
func (m *Manager) EnsureRepository(ctx context.Context, name string) error {
	repo, err := m.GetRepository(name)
	if err != nil {
		return err
	}

	err = m.EnsureFork(ctx, repo)
	if err != nil {
		return err
	}

	err = m.EnsureClone(ctx, repo)
	if err != nil {
		return err
	}

	return nil
}

// EnsureAll ensures all cached repositories.
func (m *Manager) EnsureAll(ctx context.Context) error {
	for _, repo := range m.repoCache {
		err := m.EnsureFork(ctx, repo)
		if err != nil {
			return err
		}

		err = m.EnsureClone(ctx, repo)
		if err != nil {
			return err
		}
	}
	return nil
}

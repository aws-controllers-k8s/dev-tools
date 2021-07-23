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
	"path/filepath"
	"strings"
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
	repo, err := m.getRepository(name)
	if err == nil {
		return repo, nil
	}

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

	return m.AddRepository(name, t)
}

// AddRepository creates a new Repository object and adds it to the cache.
func (m *Manager) AddRepository(name string, t RepositoryType) (*Repository, error) {
	repo := NewRepository(name, t)

	// set expected fork name
	expectedForkName := repo.Name
	if m.cfg.Github.ForkPrefix != "" {
		expectedForkName = fmt.Sprintf("%s%s", m.cfg.Github.ForkPrefix, repo.Name)
	}

	var gitHead string
	var gitRepo *git.Repository
	fullPath := filepath.Join(m.cfg.RootDirectory, repo.Name)

	gitRepo, err := m.git.Open(fullPath)
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

	repo.gitRepo = gitRepo
	repo.GitHead = gitHead
	repo.FullPath = fullPath
	repo.ExpectedForkName = expectedForkName
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
		serviceRepoName := fmt.Sprintf("%s-controller", serviceName)
		_, err := m.LoadRepository(serviceRepoName, RepositoryTypeController)
		if err != nil {
			return err
		}
	}
	return nil
}

// getRepository returns a repository from the cache
func (m *Manager) getRepository(repoName string) (*Repository, error) {
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
		repo, err := m.getRepository(repoName)
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
	// TODO(a-hilaly) ideally we need to store all repository names (service name, fork name,
	// local clone name etc...)
	repoName = strings.TrimSuffix(repoName, "-controller")
	repo, err := m.getRepository(repoName)
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %v", repoName, err)
	}
	if repo.gitRepo != nil {
		return ErrRepositoryAlreadyExist
	}

	// clone fork repository with original name
	err = m.git.Clone(
		ctx,
		m.urlBuilder(m.cfg.Github.Username, repo.ExpectedForkName),
		repo.FullPath,
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

	fork, err := m.ghc.GetUserRepositoryFork(ctx, m.cfg.Github.Username, repo.Name)
	if err == nil {
		if *fork.Name != repo.ExpectedForkName {
			err = m.ghc.RenameRepository(ctx, m.cfg.Github.Username, *fork.Name, repo.ExpectedForkName)
			if err != nil {
				return err
			}
		}
	} else if err == github.ErrForkNotFound {
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

func (m *Manager) EnsureClone(ctx context.Context, repo *Repository) error {
	err := m.clone(ctx, repo.Name)
	if err != nil && err != ErrRepositoryAlreadyExist {
		return err
	}

	return nil
}

// EnsureRepository ensures the current user owns a fork of the given repository
// and has cloned it.
func (m *Manager) EnsureRepository(ctx context.Context, name string) error {
	repo, err := m.getRepository(name)
	if err != nil && err != ErrRepositoryDoesntExist {
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

	err = m.EnsureRemotes(ctx, repo)
	if err != nil {
		return err
	}

	return nil
}

// EnsureRemotes ensures that the local repositories have both origin and upstream
// remotes setup and point to the correct URLs.
func (m *Manager) EnsureRemotes(ctx context.Context, repo *Repository) error {
	remotes, err := util.GetRepositoryRemotes(repo.gitRepo)
	if err != nil {
		return err
	}

	expecetedOriginURL := m.urlBuilder(m.cfg.Github.Username, repo.ExpectedForkName)
	// First check that one fo the  origin URLs points to the fork url
	originURLs, ok := remotes[originRemoteName]
	if !ok || !util.InStrings(expecetedOriginURL, originURLs) {
		originURLs = append(originURLs, expecetedOriginURL)
		err = util.UpdateRepositoryRemotes(repo.gitRepo, originRemoteName, originURLs)
		if err != nil {
			return fmt.Errorf("error updating origin URL: %v", err)
		}
	}

	expectedUpstreamURL := m.urlBuilder(github.ACKOrg, repo.Name)
	// Then check that one of the upstream URLs points to the original
	// repository
	upstreamURLs, ok := remotes[upstreamRemoteName]
	if !ok || !util.InStrings(expectedUpstreamURL, upstreamURLs) {
		upstreamURLs = append(upstreamURLs, expectedUpstreamURL)
		err = util.UpdateRepositoryRemotes(repo.gitRepo, upstreamRemoteName, upstreamURLs)
		if err != nil {
			return fmt.Errorf("error updating origin URL: %v", err)
		}
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

		err = m.EnsureRemotes(ctx, repo)
		if err != nil {
			return err
		}
	}
	return nil
}

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

package github

import (
	"context"
	"errors"
	"time"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

var _ RepositoryService = &Client{}

var ErrForkNotFound = errors.New("fork not found")

const (
	ACKOrg                = "aws-controllers-k8s"
	defaultRequestTimeout = 10 * time.Second
)

// NewClient takes a token and instantiate a new Client object
func NewClient(token string) *Client {
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oc := oauth2.NewClient(ctx, ts)
	return &Client{github.NewClient(oc)}
}

// RepositoryService is the interface implemented by the Github client wrapper. It exposes
// functionalities to simplify the interactions with the repository endpoint of Github APIv3
type RepositoryService interface {
	ForkRepository(ctx context.Context, repoName string) error
	RenameRepository(ctx context.Context, owner, name, newName string) error
	GetRepository(ctx context.Context, owner, repoName string) (*github.Repository, error)
	ListRepositoryForks(ctx context.Context, repoName string) ([]*github.Repository, error)
	GetUserRepositoryFork(ctx context.Context, owner, repoName string) (*github.Repository, error)
}

// Client is a github.Client wrapper
type Client struct {
	*github.Client
}

// ForkRepository forks a Github repository from the ACK organisation.
func (c *Client) ForkRepository(ctx context.Context, repoName string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	opt := &github.RepositoryCreateForkOptions{}
	_, _, err := c.Client.Repositories.CreateFork(ctx, ACKOrg, repoName, opt)
	if err != nil {
		// AcceptedError occurs when GitHub returns 202 Accepted response with an
		// empty body, which means a job was scheduled on the GitHub side to process
		// the information needed and cache it.
		// https://github.com/google/go-github/blob/master/github/github.go#L699-L704
		if _, ok := err.(*github.AcceptedError); ok {
			return nil
		}
		return err
	}
	return nil
}

// RenameRepository renames a Github repository. The request should have admin access on the
// target repositories to be able to rename it.
func (c *Client) RenameRepository(ctx context.Context, owner, name, newName string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	opt := &github.Repository{
		Name: &newName,
	}
	_, _, err := c.Client.Repositories.Edit(ctx, owner, name, opt)
	if err != nil {
		return err
	}
	return nil
}

// GetRepository takes an owner and repoName and returns the Github repository informations
func (c *Client) GetRepository(ctx context.Context, owner, repoName string) (*github.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	repo, _, err := c.Client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// ListRepositoryForks list the forks of a given repository in the ACK organisation. It returns
// a list fork information which includes the owner and the fork name (forkInfo).
func (c *Client) ListRepositoryForks(ctx context.Context, repoName string) ([]*github.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	var forks []*github.Repository
	var err error
	var repos []*github.Repository
	var resp *github.Response = &github.Response{
		// FirstPage is always of index 1
		NextPage: 1,
	}

	// iterate over all the pages
	for resp.NextPage != 0 {
		opt := &github.RepositoryListForksOptions{
			ListOptions: github.ListOptions{
				Page: resp.NextPage,
				// Fetch the maximum possible the make smallest number of
				// possible requests
				PerPage: 100,
			},
		}

		repos, resp, err = c.Client.Repositories.ListForks(ctx, ACKOrg, repoName, opt)
		if err != nil {
			return nil, err
		}

		forks = append(forks, repos...)
	}

	return forks, nil
}

// GetUserRepositoryFork takes an ACK repository name and tries to find it fork in the user public repositories.
func (c *Client) GetUserRepositoryFork(ctx context.Context, owner string, repoName string) (*github.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	repos, err := c.ListRepositoryForks(ctx, repoName)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if *repo.Owner.Login == owner {
			return repo, nil
		}
	}
	return nil, ErrForkNotFound
}

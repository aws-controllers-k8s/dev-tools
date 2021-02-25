package github

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var ErrorForkNotFound = errors.New("fork not found")

const (
	ACKOrg                = "aws-controllers-k8s"
	defaultRequestTimeout = 10 * time.Second
)

type forkInfo struct {
	Name  string
	Owner string
}

// NewClient takes a token and instantiate a new Client object
func NewClient(token string) *Client {
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oc := oauth2.NewClient(ctx, ts)
	return &Client{github.NewClient(oc)}
}

// Client is a wrapper arround github.Client
type Client struct {
	*github.Client
}

// ForkRepository forks a github repository from the ACK organisation
func (c *Client) ForkRepository(ctx context.Context, repoName string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	opt := &github.RepositoryCreateForkOptions{}
	_, _, err := c.Client.Repositories.CreateFork(ctx, ACKOrg, repoName, opt)
	if err != nil {
		// https://github.com/google/go-github/blob/master/github/github.go#L699-L704
		if _, ok := err.(*github.AcceptedError); ok {
			return nil
		}
		return err
	}
	return nil
}

// RenameRepository renames a repository. The request should have admin access on the
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

// GetRepository fetches a repository with a given owner and repoName
func (c *Client) GetRepository(ctx context.Context, owner, repoName string) (*github.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()
	repo, _, err := c.Client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// ListRepositoryForks list the forks of a given repository in the ACK organisation. it returns
// a list fork information which includes the owner and the fork name.
func (c *Client) ListRepositoryForks(ctx context.Context, repoName string) ([]*forkInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	var forks []*forkInfo

	var err error
	var repos []*github.Repository
	var resp *github.Response = &github.Response{
		// FirstPage is always of index 1
		NextPage: 1,
	}

	// loop over all the pages
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

		for _, repo := range repos {
			forks = append(forks, &forkInfo{
				Name:  *repo.Name,
				Owner: *repo.Owner.Login,
			})
		}
	}

	return forks, nil
}

// GetRepositoryForkForUser tries to find a fork repository (from ACK organisation)
// for a given user.
func (c *Client) GetRepositoryForkForUser(ctx context.Context, repoName, username string) (*forkInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	forks, err := c.ListRepositoryForks(ctx, repoName)
	if err != nil {
		return nil, err
	}

	for _, fork := range forks {
		if strings.ToLower(fork.Owner) == strings.ToLower(username) {
			return fork, nil
		}
	}

	return nil, ErrorForkNotFound
}

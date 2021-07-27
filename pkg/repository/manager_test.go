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
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	gogithub "github.com/google/go-github/v35/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	ackdevgit "github.com/aws-controllers-k8s/dev-tools/pkg/git"
	"github.com/aws-controllers-k8s/dev-tools/pkg/github"
	"github.com/aws-controllers-k8s/dev-tools/pkg/testutil"

	"github.com/aws-controllers-k8s/dev-tools/mocks"
)

var (
	testingCtx = context.TODO()
)

func stringPtr(s string) *string { return &s }

func TestManager_LoadRepository(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	testRepo, err := testutil.NewInMemoryGitRepository()
	require.NoError(err)

	fakeGit := &mocks.OpenCloner{}
	fakeGit.On("Open", "runtime").Return(testRepo, nil)
	fakeGit.On("Open", "s3-controller").Return(nil, git.ErrRepositoryNotExists)
	fakeGit.On("Open", "sqs-controller").Return(nil, ErrUnconfiguredRepository)

	type fields struct {
		cfg       *config.Config
		git       ackdevgit.OpenCloner
		repoCache map[string]*Repository
	}
	type args struct {
		name string
		t    RepositoryType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Repository
		wantErr bool
	}{
		{
			name: "repository exists",
			fields: fields{
				cfg:       testutil.NewConfig(),
				git:       fakeGit,
				repoCache: make(map[string]*Repository),
			},
			args: args{
				name: "runtime",
				t:    RepositoryTypeCore,
			},
			wantErr: false,
		},
		{
			name: "repository doesn't exists",
			fields: fields{
				cfg:       testutil.NewConfig(),
				git:       fakeGit,
				repoCache: make(map[string]*Repository),
			},
			args: args{
				name: "s3",
				t:    RepositoryTypeController,
			},
			wantErr: true,
		},
		{
			name: "unconfigured repository",
			fields: fields{
				cfg:       testutil.NewConfig(),
				git:       fakeGit,
				repoCache: make(map[string]*Repository),
			},
			args: args{
				name: "sqs",
				t:    RepositoryTypeController,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				cfg:       tt.fields.cfg,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			_, err := m.LoadRepository(tt.args.name, tt.args.t)
			if (err != nil) != tt.wantErr {
				assert.Fail(fmt.Sprintf("Manager.LoadRepository() error = %v, wantErr %v", err, tt.wantErr))
			}
		})
	}
}

func TestManager_LoadAll(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testRepo, err := testutil.NewInMemoryGitRepository()
	require.NoError(err)

	fakeGit := &mocks.OpenCloner{}
	fakeGit.On("Open", "runtime").Return(testRepo, nil)
	fakeGit.On("Open", "code-generator").Return(testRepo, nil)
	fakeGit.On("Open", "s3-controller").Return(nil, git.ErrRepositoryNotExists)
	fakeGit.On("Open", "ecr-controller").Return(nil, bytes.ErrTooLarge)
	fakeGit.On("Open", "sagemaker-controller").Return(nil, bytes.ErrTooLarge)

	type fields struct {
		cfg       *config.Config
		git       ackdevgit.OpenCloner
		repoCache map[string]*Repository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all repositories exists",
			fields: fields{
				cfg:       testutil.NewConfig(),
				git:       fakeGit,
				repoCache: make(map[string]*Repository),
			},
			wantErr: false,
		},
		{
			name: "repository not found",
			fields: fields{
				cfg:       testutil.NewConfig("s3"),
				git:       fakeGit,
				repoCache: make(map[string]*Repository),
			},
			// `ackdev list repositories` should not hide non cloned repositories.
			// It should just show that there are no active branches locally.
			wantErr: false,
		},
		{
			name: "unexpected repository error",
			fields: fields{
				cfg:       testutil.NewConfig("ecr", "sagemaker"),
				git:       fakeGit,
				repoCache: make(map[string]*Repository),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				cfg:       tt.fields.cfg,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			err := m.LoadAll()
			if (err != nil) != tt.wantErr {
				assert.Fail(fmt.Sprintf("Manager.LoadAll() error = %v, wantErr %v", err, tt.wantErr))
			}
		})
	}
}

func TestManager_clone(t *testing.T) {
	//assert := assert.New(t)
	require := require.New(t)

	testRepo, err := testutil.NewInMemoryGitRepository()
	require.NoError(err)

	fakeGit := &mocks.OpenCloner{}
	fakeGit.On("Open", "s3-controller").Return(testRepo, nil)
	fakeGit.On("Open", "mq-controller").Return(nil, git.ErrRepositoryNotExists)
	fakeGit.On("Open", "ecr-controller").Return(nil, git.ErrRepositoryNotExists)
	fakeGit.On("Open", "sagemaker-controller").Return(nil, git.ErrRepositoryNotExists).Once()
	fakeGit.On("Open", "sagemaker-controller").Return(testRepo, nil).Once()

	fakeGit.On(
		"Clone",
		testingCtx,
		"https://github.com/ack-bot/ack-ecr-controller.git",
		"ecr-controller",
	).Return(transport.ErrAuthenticationRequired)
	fakeGit.On(
		"Clone",
		testingCtx,
		"https://github.com/ack-bot/ack-mq-controller.git",
		"mq-controller",
	).Return(gitconfig.ErrRemoteConfigNotFound)
	fakeGit.On(
		"Clone",
		testingCtx,
		"https://github.com/ack-bot/ack-sagemaker-controller.git",
		"sagemaker-controller",
	).Return(nil)

	type fields struct {
		cfg        *config.Config
		ghc        github.RepositoryService
		git        ackdevgit.OpenCloner
		urlBuilder func(string, string) string
		repoCache  map[string]*Repository
	}
	type args struct {
		repoName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "repository not configured",
			fields: fields{
				cfg:       testutil.NewConfig(),
				git:       fakeGit,
				repoCache: make(map[string]*Repository),
			},
			args: args{
				repoName: "dynamodb",
			},
			wantErr: true,
		},
		{
			name: "repository already exists",
			fields: fields{
				cfg:       testutil.NewConfig("s3", "elasticache"),
				git:       fakeGit,
				repoCache: make(map[string]*Repository),
			},
			args: args{
				repoName: "s3",
			},
			wantErr: true,
		},
		{
			name: "unauthenticated git",
			fields: fields{
				cfg:        testutil.NewConfig("s3", "elasticache"),
				git:        fakeGit,
				urlBuilder: httpsRemoteURL,
				repoCache:  make(map[string]*Repository),
			},
			args: args{
				repoName: "ecr",
			},
			wantErr: true,
		},
		{
			name: "cloning error",
			fields: fields{
				cfg:        testutil.NewConfig("mq"),
				git:        fakeGit,
				urlBuilder: httpsRemoteURL,
				repoCache:  make(map[string]*Repository),
			},
			args: args{
				repoName: "mq",
			},
			wantErr: true,
		},
		{
			name: "cloning successful",
			fields: fields{
				cfg:        testutil.NewConfig("sagemaker"),
				git:        fakeGit,
				urlBuilder: httpsRemoteURL,
				repoCache:  make(map[string]*Repository),
			},
			args: args{
				repoName: "sagemaker",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				cfg:        tt.fields.cfg,
				ghc:        tt.fields.ghc,
				git:        tt.fields.git,
				urlBuilder: tt.fields.urlBuilder,
				repoCache:  tt.fields.repoCache,
			}
			_, _ = m.LoadRepository(tt.args.repoName, RepositoryTypeController)
			if err := m.clone(testingCtx, tt.args.repoName); (err != nil) != tt.wantErr {
				t.Errorf("Manager.clone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_EnsureFork(t *testing.T) {
	fakeGithubClient := &mocks.RepositoryService{}
	// s3 case
	fakeGithubClient.On(
		"GetUserRepositoryFork",
		testingCtx,
		"ack-bot",
		"s3-controller",
	).Return(nil, github.ErrForkNotFound)
	fakeGithubClient.On(
		"ForkRepository",
		testingCtx,
		"s3-controller",
	).Return(errors.New("unknown error"))

	// sagemaker case
	fakeGithubClient.On(
		"GetUserRepositoryFork",
		testingCtx,
		"ack-bot",
		"sagemaker-controller",
	).Return(&gogithub.Repository{Name: stringPtr("sagemaker-controller")}, nil)
	fakeGithubClient.On(
		"RenameRepository",
		testingCtx,
		"ack-bot",
		"sagemaker-controller",
		"ack-sagemaker-controller",
	).Return(nil)

	// ecr case
	fakeGithubClient.On(
		"GetUserRepositoryFork",
		testingCtx,
		"ack-bot",
		"ecr-controller",
	).Return(nil, github.ErrForkNotFound)
	fakeGithubClient.On(
		"ForkRepository",
		testingCtx,
		"ecr-controller",
	).Return(nil)
	fakeGithubClient.On(
		"RenameRepository",
		testingCtx,
		"ack-bot",
		"ecr-controller",
		"ack-ecr-controller",
	).Return(nil)

	type fields struct {
		cfg       *config.Config
		ghc       github.RepositoryService
		git       ackdevgit.OpenCloner
		repoCache map[string]*Repository
	}
	type args struct {
		repo *Repository
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fork error",
			fields: fields{
				cfg:       testutil.NewConfig("s3"),
				ghc:       fakeGithubClient,
				repoCache: make(map[string]*Repository),
			},
			args: args{
				repo: &Repository{
					Name:             "s3-controller",
					ExpectedForkName: "s3-sagemaker-controller",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure fork successful - rename",
			fields: fields{
				cfg:       testutil.NewConfig("sagemaker"),
				ghc:       fakeGithubClient,
				repoCache: make(map[string]*Repository),
			},
			args: args{
				repo: &Repository{
					Name:             "sagemaker-controller",
					ExpectedForkName: "ack-sagemaker-controller",
				},
			},
			wantErr: false,
		},
		{
			name: "ensure fork successful - fork and rename",
			fields: fields{
				cfg:       testutil.NewConfig("ecr"),
				ghc:       fakeGithubClient,
				repoCache: make(map[string]*Repository),
			},
			args: args{
				repo: &Repository{
					Name:             "ecr-controller",
					ExpectedForkName: "ack-ecr-controller",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				cfg:       tt.fields.cfg,
				ghc:       tt.fields.ghc,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			if err := m.EnsureFork(testingCtx, tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ensureFork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

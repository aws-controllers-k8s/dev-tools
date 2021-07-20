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
	// this field might be nil, if the repository doesn't exist locally
	gitRepo *git.Repository

	// Name of the ACK upstream repo
	Name string
	// Repository Type
	Type RepositoryType
	// Expected fork name. Generally looking like ack-sagemaker
	ExpectedForkName string
	// Expected local full path
	FullPath string
	// Git HEAD commit or current branch
	GitHead string
}

func httpsRemoteURL(owner, name string) string {
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, name)
}

func sshRemoteURL(owner, name string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", owner, name)
}

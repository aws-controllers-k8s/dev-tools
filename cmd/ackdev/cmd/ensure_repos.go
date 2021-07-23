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

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

var ensureRepositoriesCmd = &cobra.Command{
	Use:     "repo",
	Aliases: []string{"repos", "repositories", "repository"},
	RunE:    ensureAllRepositories,
	Args:    cobra.NoArgs,
	Short:   "Ensure repositories are forked and cloned locally",
}

func ensureAllRepositories(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	err = repoManager.LoadAll()
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	return repoManager.EnsureAll(ctx)
}

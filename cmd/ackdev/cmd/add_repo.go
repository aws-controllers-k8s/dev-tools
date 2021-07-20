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
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
)

var (
	optAddRepoType string
)

func init() {
	addRepositoryCmd.PersistentFlags().StringVarP(&optAddRepoType, "type", "t", "controller", "repository type")
}

var addRepositoryCmd = &cobra.Command{
	Use:     "repository <service> ...",
	Aliases: []string{"repo", "repos", "repository", "repositories"},
	RunE:    addRepository,
	Args:    cobra.MinimumNArgs(1),
}

func addRepository(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	for _, service := range args {
		service = strings.ToLower(service)

		// Check it doesn't already exist in the configuration
		if util.InStrings(service, cfg.Repositories.Services) {
			fmt.Printf("repository for service %s has already been added\n", service)
			continue
		}

		_, err := repoManager.AddRepository(service, repository.GetRepositoryTypeFromString(optAddRepoType))
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		if err := repoManager.EnsureRepository(ctx, service); err != nil {
			return err
		}

		cfg.Repositories.Services = append(cfg.Repositories.Services, service)
		if err := config.Save(cfg, ackConfigPath); err != nil {
			return err
		}
	}

	return nil
}

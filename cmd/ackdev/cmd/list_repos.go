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

var (
	listTableHeaderColumns = []string{"Name", "Type"}

	optListFilterExpression string
	optListShowBranch       bool
)

func init() {
	listRepositoriesCmd.PersistentFlags().StringVarP(&optListFilterExpression, "filter", "f", "", "filter expression")
	listRepositoriesCmd.PersistentFlags().BoolVar(&optListShowBranch, "show-branch", true, "display project current branch or not")
}

var listRepositoriesCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo", "repos", "repositories"},
	RunE:    printRepositories,
	Args:    cobra.NoArgs,
}

func printRepositories(cmd *cobra.Command, args []string) error {
	filters, err := repository.BuildFilters(optListFilterExpression)
	if err != nil {
		return err
	}

	repos, err := listRepositories(filters...)
	if err != nil {
		return err
	}

	tablePrintRepositories(repos)
	return nil
}

func listRepositories(filters ...repository.Filter) ([]*repository.Repository, error) {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return nil, err
	}
	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return nil, err
	}

	// Try to load all repositories
	err = repoManager.LoadAll()
	if err != nil {
		return nil, err
	}

	// List repositories
	//TODO(hilalymh) add sort-by flag/option
	repos := repoManager.List(filters...)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func tablePrintRepositories(repos []*repository.Repository) {
	tableHeaderColumns := listTableHeaderColumns
	if optListShowBranch {
		tableHeaderColumns = append(tableHeaderColumns, "Branch")
	}

	tw := newTable()
	defer tw.Render()

	tw.SetHeader(tableHeaderColumns)

	for _, repo := range repos {
		rawArgs := []string{repo.Name, repo.Type.String()}
		if optListShowBranch {
			rawArgs = append(rawArgs, repo.GitHead)
		}
		tw.Append(rawArgs)
	}
}

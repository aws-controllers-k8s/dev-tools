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

import "github.com/spf13/cobra"

var (
	optListOutputFormat string
)

func init() {
	listCmd.AddCommand(listDependenciesCmd)
	listCmd.AddCommand(listRepositoriesCmd)
	listCmd.AddCommand(getConfigCmd)

	getConfigCmd.PersistentFlags().StringVarP(&optListOutputFormat, "output", "o", "yaml", "output format (json|yaml)")
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"get", "ls"},
	Args:    cobra.NoArgs,
	Short:   "Display one or many resources",
}

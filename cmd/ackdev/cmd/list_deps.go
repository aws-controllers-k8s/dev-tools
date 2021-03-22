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

	"github.com/aws-controllers-k8s/dev-tools/pkg/deps"
)

var (
	listDepsTableHeaderColumns = []string{"Name", "Status"}
	optDepsListShowPath        bool
	optDepsListShowVersion     bool
)

func init() {
	listDependenciesCmd.PersistentFlags().BoolVar(&optDepsListShowPath, "show-path", true, "display binary path")
	listDependenciesCmd.PersistentFlags().BoolVar(&optDepsListShowVersion, "show-version", true, "display binary version")
}

var listDependenciesCmd = &cobra.Command{
	Use:     "dependency",
	Aliases: []string{"dep", "deps", "dependencies"},
	RunE:    printDependencies,
}

type depRecord struct {
	Name    string
	Version string
	Path    string
	Status  string
}

func printDependencies(cmd *cobra.Command, args []string) error {
	dependencies, err := listDependencies()
	if err != nil {
		return err
	}

	tablePrintDependencies(dependencies)
	return nil
}

// tablePrintDependencies prints the ACK development dependencies in table
func tablePrintDependencies(dependencies []*depRecord) error {
	// table headers
	tableHeaderColumns := listDepsTableHeaderColumns
	if optDepsListShowVersion {
		tableHeaderColumns = append(tableHeaderColumns, "Version")
	}
	if optDepsListShowPath {
		tableHeaderColumns = append(tableHeaderColumns, "Path")
	}

	tw := newTable()
	defer tw.Render()

	// Add header
	tw.SetHeader(tableHeaderColumns)

	for _, tool := range dependencies {
		rawArgs := []string{tool.Name, tool.Status}
		if optDepsListShowVersion {
			rawArgs = append(rawArgs, tool.Version)
		}
		if optDepsListShowPath {
			rawArgs = append(rawArgs, tool.Path)
		}
		tw.Append(rawArgs)
	}

	return nil
}

// listDependencies returns the list of ACK development dependencies
// along with their versions and binary paths.
func listDependencies() ([]*depRecord, error) {
	list := make([]*depRecord, 0, len(deps.DevelopmentTools))
	for _, tool := range deps.DevelopmentTools {
		status := ""
		path, err := tool.BinPath()
		if err != nil {
			status = "NOT FOUND"
		} else {
			status = "OK"
		}

		version := "-"
		if status == "OK" {
			v, err := tool.Version()
			if err != nil && err != deps.ErrorVersionNotFound {
				return nil, err
			}
			version = v
		}
		list = append(list, &depRecord{
			Name:    tool.BinaryName,
			Version: version,
			Path:    path,
			Status:  status,
		})
	}
	return list, nil
}

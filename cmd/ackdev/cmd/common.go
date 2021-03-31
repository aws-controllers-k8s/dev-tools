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
	"go/build"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
)

const (
	ackdevConfigFileName = ".ackdev.yaml"
)

var (
	homeDirectory        string
	defaultConfigPath    string
	goPath               = build.Default.GOPATH
	defaultRootDirectory = filepath.Join(goPath, "src/github.com/aws-controllers-k8s")
)

func init() {
	// Set homeDirectory and defaultConfigPath
	hd, err := homedir.Dir()
	if err != nil {
		fmt.Printf("unable to determine $HOME: %s\n", err)
		os.Exit(1)
	}
	homeDirectory = hd
	defaultConfigPath = filepath.Join(homeDirectory, ackdevConfigFileName)
}

func newTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)

	// Kubectl tables like style
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(true)
	return table
}

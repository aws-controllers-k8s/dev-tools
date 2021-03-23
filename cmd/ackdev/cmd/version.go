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
	"runtime"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/version"
)

const versionTmpl = `Date: %s
Build: %s
Version: %s
Git Hash: %s
`

var versionCmd = &cobra.Command{
	Use:   "version",
	Args:  cobra.NoArgs,
	RunE:  printVersion,
	Short: "Print ackdev binary version informations",
}

func printVersion(*cobra.Command, []string) error {
	goVersion := fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf(versionTmpl, version.BuildDate, goVersion, version.GitVersion, version.GitCommit)
	return nil
}

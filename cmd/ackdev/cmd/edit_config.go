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
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var editConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify ackdev configuration file",
	Args:  cobra.NoArgs,
	RunE:  editConfig,
}

// editConfig opens ackdev configuration file in an editor. By default
// opens the configuration file using vi.
func editConfig(cmd *cobra.Command, args []string) error {
	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	c := exec.Command(executable, ackConfigPath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

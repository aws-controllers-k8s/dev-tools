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

	"github.com/spf13/cobra"
)

const (
	defaultEditor = "vi"
)

var (
	editor string
)

func init() {
	editor = os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}

	editCmd.AddCommand(editConfigCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Args:  cobra.NoArgs,
	Short: "Edit a resource from the default editor.",
}

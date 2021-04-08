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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
)

var (
	optSetupRootDirectory   string
	optSetupInitialServices string
)

func init() {
	setupCmd.PersistentFlags().StringVar(&optSetupRootDirectory, "root-directory", defaultRootDirectory, "root directory for ACK repositories")
	setupCmd.PersistentFlags().StringVarP(&optSetupInitialServices, "services", "s", "", "services injected in the generated configuration file")
}

var setupCmd = &cobra.Command{
	Use:     "setup",
	RunE:    setupACKDev,
	Args:    cobra.NoArgs,
	Short:   "Generate ackdev configuration file",
	Example: "ackdev setup --root-dir=. --services=s3,ecr,sqs",
}

func setupACKDev(cmd *cobra.Command, args []string) error {
	_, err := config.Load(ackConfigPath)
	if err == nil {
		return fmt.Errorf("ackdev is already setup")
	}

	initialServices := strings.Split(optSetupInitialServices, ",")
	rootDir, err := filepath.Abs(optSetupRootDirectory)
	if err != nil {
		return err
	}

	// Ensure that the root directory exists
	err = os.MkdirAll(rootDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	newConfig := config.Config{
		RootDirectory: rootDir,
		Repositories: config.RepositoriesConfig{
			Services: initialServices,
			Core:     config.DefaultConfig.Repositories.Core,
		},
	}

	err = config.Save(&newConfig, ackConfigPath)
	if err != nil {
		return err
	}
	return nil
}

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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/soak"
)

var (
	defaultSoakClusterConfig = filepath.Join(defaultRootDirectory, "test-infra/soak/cluster-config.yaml")
)

var (
	optSoakBootstrapClusterConfig string
)

func init() {
	soakBootstrapCmd.PersistentFlags().StringVar(&optSoakBootstrapClusterConfig, "cluster-config", defaultSoakClusterConfig, "eksctl cluster config file")
}

var soakBootstrapCmd = &cobra.Command{
	Use:     "bootstrap",
	RunE:    soakBootstrap,
	Args:    cobra.NoArgs,
	Short:   "Bootstrap the soak test cluster",
	Example: "ackdev soak bootstrap --cluster-config=./cluster-config.yml --service=s3",
}

func soakBootstrap(cmd *cobra.Command, args []string) error {
	_, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	clusterConfig, err := filepath.Abs(optSoakBootstrapClusterConfig)
	if err != nil {
		return err
	}

	if optSoakService == "" {
		return fmt.Errorf("you must specify a service to bootstrap")
	}

	fmt.Printf("Bootstrapping ECR Public Repo ... ")

	repoUri, err := soak.EnsureECRRepository(optSoakService)
	if err != nil {
		return err
	}

	fmt.Printf("👍 (%s)\n", repoUri)

	fmt.Println("Bootstrapping EKS Cluster (this may take >30 minutes) ... ")
	err = soak.EnsureCluster(clusterConfig, optSoakService)
	if err != nil {
		return err
	}
	fmt.Println("👍")

	return nil
}

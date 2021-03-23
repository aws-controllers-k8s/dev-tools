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
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
)

var getConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Display ackdev configuration file",
	Args:  cobra.NoArgs,
	RunE:  printConfig,
}

func printConfig(*cobra.Command, []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	var b []byte
	switch optListOutputFormat {
	case "json":
		b, err = json.Marshal(cfg)
		if err != nil {
			return err
		}
	case "yaml":
		b, err = yaml.Marshal(cfg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported output type: %s", optListOutputFormat)
	}

	fmt.Println(string(b))
	return nil
}

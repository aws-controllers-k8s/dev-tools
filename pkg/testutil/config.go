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

package testutil

import "github.com/aws-controllers-k8s/dev-tools/pkg/config"

// Returns a new config.Config object used for testing purposes.
func NewConfig(services ...string) *config.Config {
	return &config.Config{
		Repositories: config.RepositoriesConfig{
			Core: []string{
				"runtime",
				"code-generator",
			},
			Services: services,
		},
		Github: config.GithubConfig{
			ForkPrefix: "ack-",
			Username:   "ack-bot",
		},
	}
}

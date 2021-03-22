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

package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Config is the ackdev global configuration. It contains information and default values
// used by ackdev to manage local repositories, forks, dependencies, controllers...
type Config struct {
	// RootDirectory is the parent directory of the all ACK local repositories.
	// If it's not specified ackdev will use $GOPATH/src/github.com/aws-controllers-k8s
	RootDirectory string `yaml:"rootDirectory" json:"rootDirectory"`
	// Git contains information used by ackdev to manage local git repositories.
	Git GitConfig `yaml:"git" json:"git"`
	// Github contains information used by ackdev to manage Github forks.
	Github GithubConfig `yaml:"github" json:"github"`
	// Repositories let you specify the Github/Local repositories ackdev should
	// manage for you. By default ackdev manage all the core repositories (community,
	// code-generator, test-infra, dev-tools and runtime).
	Repositories RepositoriesConfig `yaml:"repositories" json:"repositories"`
	// RunConfig let specify the arguments and flags used to run a controller locally,
	// without having to build it image or deploy it into a cluster.
	RunConfig RunConfig `yaml:"run" json:"run"`
}

// RepositoriesConfig represent repositories that are be managed by ackdev.
type RepositoriesConfig struct {
	// Core is the list of ACK core repositories. The default configuration contains:
	// commmunity, code-generator, test-infra, dev-tools and runtime.
	Core []string `yaml:"core" json:"core"`
	// Services is the list of service controllers managed by ackdev.
	Services []string `yaml:"services" json:"services"`
}

// GithubConfig represents the Github information needed to personal forks.
type GithubConfig struct {
	// Token is the token used to make Github API calls. This token needs at least
	// the 'repo' scope. To generate this token please follow instructions in:
	// https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token
	Token string `yaml:"token" json:"token"`
	// Username is the ackdev contributor Github username.
	Username string `yaml:"username" json:"username"`
	// ForkPrefix is the prefix prepended to the personal forks of ACK repositories.
	// For example if ForkPrefix is 'ack-', ackdev will fork code-generator repository
	// and rename to 'ack-code-generator.
	ForkPrefix string `yaml:"forkPrefix" json:"forkPrefix"`
}

// Git contains information used by ackdev to manage local git repositories.
type GitConfig struct {
	// SSHKeyPath is the full path the SSH key used to clone Github repositories.
	SSHKeyPath string `yaml:"sshKeyPath" json:"sshKeyPath"`
}

// RunConfig contains flags and arguments passed to service controllers binaries when
// they are executed locally.
type RunConfig struct {
	// Flags is the map of flags/values passed to the controller binaries. For example
	// to pass --aws-region=us-west-1 you'll need to set Flags to {"aws-region","us-west-1"}
	Flags map[string]string `yaml:"flags" json:"flags"`
}

// DefaultConfig is the default configuration used to generated ackdev config
var DefaultConfig = Config{
	Repositories: RepositoriesConfig{
		Core: []string{
			"runtime",
			"dev-tools",
			"community",
			"code-generator",
			"test-infra",
		},
	},
	Github: GithubConfig{
		ForkPrefix: "ack-",
	},
}

// Load reads a local configuration file and returns an ackdev configuration object.
func Load(configPath string) (*Config, error) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save serialise a configuration object and writes it to given filepath.
func Save(cfg *Config, filename string) error {
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return nil
	}
	err = ioutil.WriteFile(filename, bytes, 0777)
	if err != nil {
		return err
	}
	return nil
}

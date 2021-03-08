package config

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RootDirectory string             `yaml:"rootDirectory" json:"rootDirectory"`
	Git           GitConfig          `yaml:"git" json:"git"`
	Github        GithubConfig       `yaml:"github" json:"github"`
	Repositories  RepositoriesConfig `yaml:"repositories" json:"repositories"`
	RunConfig     RunConfig          `yaml:"run" json:"run"`
}

type RepositoriesConfig struct {
	Core     []string `yaml:"core" json:"core"`
	Services []string `yaml:"services" json:"services"`
}

type GithubConfig struct {
	Token      string `yaml:"token" json:"token"`
	Username   string `yaml:"username" json:"username"`
	ForkPrefix string `yaml:"forkPrefix" json:"forkPrefix"`
}

type GitConfig struct {
	SSHKeyPath string `yaml:"sshKeyPath" json:"sshKeyPath"`
}

type RunConfig struct {
	Flags map[string]string `yaml:"flags" json:"flags"`
}

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

func LoadFromJSONBytes(configBytes []byte) (*Config, error) {
	cfg := DefaultConfig
	err := json.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

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

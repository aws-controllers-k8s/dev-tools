package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

func init() {}

var ensureRepositoriesCmd = &cobra.Command{
	Use:     "repo",
	Aliases: []string{"repo", "repos", "repositories"},
	RunE:    ensureAll,
	Args:    cobra.NoArgs,
}

func ensureAll(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	err = repoManager.LoadAll()
	if err != nil {
		return err
	}

	err = repoManager.EnsureAll()
	if err != nil {
		return err
	}

	return nil
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/deps"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

var ensureDependenciesCmd = &cobra.Command{
	Use:     "dependency",
	Aliases: []string{"dep", "deps", "dependencies"},
	RunE:    ensureDependencies,
}

func ensureDependencies(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}
	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}
	communityRepository, err := repoManager.LoadRepository("community", repository.RepositoryTypeCore)
	if err != nil {
		return err
	}

	for _, tool := range deps.DevelopmentTools {
		_, err := tool.BinPath()
		if err == nil {
			continue
		}

		fmt.Printf("installling %s... ", tool.Name)

		scriptPath := fmt.Sprintf("scripts/install-%s.sh", tool.Name)
		err = executeCommand(communityRepository.FullPath, "bash", []string{scriptPath}, true)
		if err != nil {
			fmt.Printf("FAIL\n")
		} else {
			fmt.Printf("OK\n")
		}
	}

	return nil
}

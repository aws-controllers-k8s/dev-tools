package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
)

var (
	optAddEnsure         bool
	optAddRepositoryType string
)

func init() {
	addCmd.PersistentFlags().StringVar(&optAddRepositoryType, "type", "controller", "repository type (controller|core)")
}

var addCmd = &cobra.Command{
	Use:   "add",
	RunE:  addRepository,
	Args:  cobra.MinimumNArgs(1),
	Short: "Add new controllers or core repositories",
}

func addRepository(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	for _, name := range args {
		switch optAddRepositoryType {
		case "controller":
			if !util.InStrings(name, cfg.Repositories.Services) {
				cfg.Repositories.Services = append(cfg.Repositories.Services, name)
			}
		case "core":
			if !util.InStrings(name, cfg.Repositories.Core) {
				cfg.Repositories.Core = append(cfg.Repositories.Core, name)
			}
		default:
			return fmt.Errorf("unsupported reposiroty type %s", optAddRepositoryType)
		}
	}

	err = config.Save(cfg, ackConfigPath)
	if err != nil {
		return err
	}

	return nil
}

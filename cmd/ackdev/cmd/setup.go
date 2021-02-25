package cmd

import (
	"log"
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
	setupCmd.PersistentFlags().StringVarP(&optSetupRootDirectory, "root-directory", "R", defaultRootDirectory, "root directory for ACK projects")
	setupCmd.PersistentFlags().StringVarP(&optSetupInitialServices, "services", "s", "", "services initialized with the ack configuration")
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	RunE:  setupACKDev,
	Args:  cobra.NoArgs,
	Short: "Generate ackdev configuration file",
}

func setupACKDev(cmd *cobra.Command, args []string) error {
	_, err := config.Load(ackConfigPath)
	if err == nil {
		log.Println("ackdev is already setup.")
		return nil
	}

	err = generateAckDevToolsConfiguration()
	if err != nil {
		return err
	}
	return nil
}

func generateAckDevToolsConfiguration() error {
	initialServices := strings.Split(optSetupInitialServices, " ")
	rootDir, err := filepath.Abs(optSetupRootDirectory)
	if err != nil {
		return err
	}

	newConfig := config.Config{
		RootDirectory: rootDir,
		Repositories: config.RepositoriesConfig{
			Services: initialServices,
			// Core repositories are inject by default
			Core: config.DefaultConfig.Repositories.Core,
		},
	}
	err = config.Save(&newConfig, ackConfigPath)
	if err != nil {
		return err
	}
	return nil
}

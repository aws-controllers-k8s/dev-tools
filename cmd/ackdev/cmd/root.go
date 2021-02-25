package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ackConfigPath string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&ackConfigPath, "config-file", defaultConfigPath, "ack config file path")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(ensureCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(generateCmd)
}

var rootCmd = &cobra.Command{
	Use:           "ackdev",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "A tool to manage ACK repositories, CRDs, development tools and testing",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

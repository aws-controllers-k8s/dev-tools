package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(setConfigCmd)
	configCmd.AddCommand(viewConfigCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Args:  cobra.NoArgs,
	Short: "View or edit ackdev configuration file",
}

package cmd

import (
	"github.com/spf13/cobra"
)

var ()

func init() {
	ensureCmd.AddCommand(ensureDependenciesCmd)
	ensureCmd.AddCommand(ensureRepositoriesCmd)
}

var ensureCmd = &cobra.Command{
	Use:   "ensure",
	Args:  cobra.NoArgs,
	Short: "Ensure repositories or dependencies",
}

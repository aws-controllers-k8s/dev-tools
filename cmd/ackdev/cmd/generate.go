package cmd

import (
	"github.com/spf13/cobra"
)

var ()

func init() {
	generateCmd.AddCommand(generateMocksCmd)
	generateCmd.AddCommand(generateClientCmd)
	generateCmd.AddCommand(generateControllerCmd)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Args:    cobra.NoArgs,
	Short:   "Generate controller code or Go client",
}

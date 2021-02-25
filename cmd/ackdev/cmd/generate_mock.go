package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateMocksCmd = &cobra.Command{
	Use:     "mocks",
	Aliases: []string{"mock"},
	Args:    cobra.NoArgs,
	RunE:    generateMocks,
}

func generateMocks(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("not implemented")
}

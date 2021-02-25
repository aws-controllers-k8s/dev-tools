package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateClientCmd = &cobra.Command{
	Use:  "client",
	Args: cobra.NoArgs,
	RunE: generateGoClientCode,
}

func generateGoClientCode(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("not implemented")
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {}

var deployControllerCmd = &cobra.Command{
	Use:  "controller",
	RunE: deployController,
	Args: cobra.NoArgs,
}

func deployController(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("not implemented")
}

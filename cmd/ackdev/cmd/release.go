package cmd

import "github.com/spf13/cobra"

var releaseCmd = &cobra.Command{
	Use:  "release",
	Args: cobra.NoArgs,
}

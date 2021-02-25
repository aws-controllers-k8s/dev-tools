package cmd

import "github.com/spf13/cobra"

func init() {
	deployCmd.AddCommand(deployCRDsCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Args:  cobra.NoArgs,
	Short: "Deploy controllers or crds to a kubernetes cluster",
}

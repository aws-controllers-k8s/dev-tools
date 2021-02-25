package cmd

import "github.com/spf13/cobra"

var (
	optListOutput string
)

func init() {
	listCmd.PersistentFlags().StringVarP(&optListOutput, "output", "o", "table", "output format (json|yaml|table)")
	listCmd.AddCommand(listRepositoriesCmd)
	listCmd.AddCommand(listDependenciesCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"get", "ls"},
	Args:    cobra.NoArgs,
	Short:   "Display one or many resources",
}

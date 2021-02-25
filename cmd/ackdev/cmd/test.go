package cmd

import (
	"github.com/spf13/cobra"
)

var (
	optTestParallel    bool
	optTestMaxWorkers  int
	optTestShowSummary bool
)

func init() {
	testCmd.AddCommand(unitTestCmd)
	testCmd.AddCommand(e2eTestCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Args:  cobra.MinimumNArgs(1),
	Short: "Run unit or end2end tests for a given or multiple controllers",
}

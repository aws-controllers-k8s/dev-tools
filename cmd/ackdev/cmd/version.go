package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/version"
)

var ()

func init() {

}

const versionTmpl = `Date: %s
Build: %s
Version: %s
Git Hash: %s
`

var versionCmd = &cobra.Command{
	Use:   "version",
	Args:  cobra.NoArgs,
	RunE:  printVersion,
	Short: "Print ackdev binary version informations",
}

func printVersion(*cobra.Command, []string) error {
	goVersion := fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf(versionTmpl, version.BuildDate, goVersion, version.GitVersion, version.GitCommit)
	return nil
}

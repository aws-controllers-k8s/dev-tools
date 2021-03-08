package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/aexec"
	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

var (
	optGenerateControllerParallel     bool
	optGenerateControllerSummaryTable bool
	optGenerateControllerStreamLogs   bool

	generateControllerResultsTableColumns = []interface{}{"SERVICE", "SUCCESS", "LOGFILE"}
)

func init() {}

var generateControllerCmd = &cobra.Command{
	Use:     "controller",
	Aliases: []string{"ctrl"},
	RunE:    generateController,
	Args:    cobra.ExactArgs(1),
}

func generateController(cmd *cobra.Command, args []string) error {
	if optGenerateControllerParallel && optGenerateControllerStreamLogs {
		return fmt.Errorf("flag conflict: --parallel, --stream-logs")
	}

	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	codeGeneratorRepo, err := repoManager.LoadRepository("code-generator", repository.RepositoryTypeCore)
	if err != nil {
		return err
	}

	return runACKGenerate(codeGeneratorRepo.FullPath, args[0])
}

func runACKGenerate(codeGeneratorPath, service string) error {
	cmd := exec.Command("make", "build-controller", "SERVICE="+service)
	cmd.Dir = codeGeneratorPath
	acmd := aexec.New(cmd, 50)
	err := acmd.Run(true)
	if err != nil {
		return err
	}
	done := make(chan struct{})

	go func() {
		for b := range acmd.StdoutStream() {
			fmt.Println(string(b))
		}
		done <- struct{}{}

	}()
	go func() {
		for b := range acmd.StderrStream() {
			fmt.Println(string(b))
		}
		done <- struct{}{}

	}()

	// wait for printers to finish
	defer func() { _, _ = <-done, <-done }()

	err = acmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

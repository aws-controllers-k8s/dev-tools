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
	unitTestsResultsColumns = []interface{}{"NAME", "TYPE"}

	optUnitSaveLogs bool
)

func init() {
}

var unitTestCmd = &cobra.Command{
	Use:  "unit",
	RunE: runUnitTests,
	Args: cobra.ExactArgs(1),
}

func runUnitTests(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	serviceRepo, err := repoManager.LoadRepository(args[0], repository.RepositoryTypeController)
	if err != nil {
		return err
	}

	return runServiceUnitTests(serviceRepo.FullPath)
}

func runServiceUnitTests(servicePath string) error {
	cmd := exec.Command("make", "test")
	cmd.Dir = servicePath
	acmd := aexec.New(cmd, 300)
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

	// wait for command to finish
	err = acmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

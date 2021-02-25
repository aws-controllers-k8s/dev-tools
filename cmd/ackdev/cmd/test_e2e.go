package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/aexec"
	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

func init() {

}

var e2eTestCmd = &cobra.Command{
	Use:  "e2e",
	RunE: runE2ETests,
	Args: cobra.MinimumNArgs(0),
}

func runE2ETests(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	services := args

	communityRepo, err := repoManager.LoadRepository("community", repository.RepositoryTypeCore)
	if err != nil {
		return err
	}

	if len(services) == 1 {
		serviceRepo, err := repoManager.LoadRepository(services[0], repository.RepositoryTypeController)
		if err != nil {
			return err
		}

		// stream e2e tests logs
		return runKindBuildTestsWithLogStream(communityRepo.FullPath, serviceRepo.Name)
	}

	return fmt.Errorf("not implemented")
}

func runKindBuildTestsWithLogStream(codeGeneratorPath, service string) error {
	cmd := exec.Command("make", "kind-test", "SERVICE="+service)
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

package cmd

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/aexec"
	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
	"github.com/aws-controllers-k8s/dev-tools/pkg/table"
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
	Args: cobra.MinimumNArgs(1),
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

	services := args

	if len(services) == 1 {
		serviceRepo, err := repoManager.LoadRepository(services[0], repository.RepositoryTypeController)
		if err != nil {
			return err
		}

		// stream unit tests logs
		return runServiceUnitTestsWithLogStream(serviceRepo.FullPath)
	}

	tw := table.NewPrinter(len(unitTestsResultsColumns))
	defer func() {
		if err := tw.Print(); err != nil {
			panic(err)
		}
	}()

	// print header
	if err := tw.AddRaw(unitTestsResultsColumns...); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// run service tests
	for _, service := range services {
		wg.Add(1)

		go func(service string) {
			defer wg.Done()
			serviceRepo, err := repoManager.LoadRepository(service, repository.RepositoryTypeController)
			if err != nil || serviceRepo.GitHead == "" {
				if err := tw.AddRaw(service, "ERROR"); err != nil {
					panic(err)
				}
				return
			}

			result := runServiceUnitTest(serviceRepo.Name, serviceRepo.FullPath)
			if result.Error != nil {
				if err := tw.AddRaw(service, "FAIL"); err != nil {
					panic(err)
				}
			} else {
				if err := tw.AddRaw(service, "PASS"); err != nil {
					panic(err)
				}
			}
		}(service)
	}

	wg.Wait()
	return nil
}

type uniTestResults struct {
	service string
	success bool
	err     error
	stdout  [][]byte
	stderr  [][]byte
}

func runServiceUnitTest(serviceName, servicePath string) aexec.CmdResult {
	cmd := exec.Command("make", "test")
	cmd.Dir = servicePath
	acmd := aexec.New(cmd, 300)
	err := acmd.Run(true)
	if err != nil {
		return aexec.CmdResult{
			Error:   err,
			Service: serviceName,
			Success: false,
		}
	}
	err = acmd.Wait()
	return aexec.CmdResult{
		Error:   err,
		Service: serviceName,
		Success: err == nil,
		Stdout:  acmd.Stdout,
		Stderr:  acmd.Stderr,
	}
}

func runServiceUnitTestsWithLogStream(servicePath string) error {
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

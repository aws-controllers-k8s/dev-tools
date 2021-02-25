package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"k8s.io/client-go/util/homedir"

	"github.com/aws-controllers-k8s/dev-tools/pkg/aexec"
)

var (
	GOPATH               = os.Getenv("GOPATH")
	GOBIN                = os.Getenv("GOBIN")
	homeDirectory        = homedir.HomeDir()
	defaultConfigPath    = filepath.Join(homeDirectory, ".ackdev.yaml")
	defaultRootDirectory = filepath.Join(GOPATH, "source/github.com/aws-controllers-k8s")
)

func executeCommand(contextDir string, command string, args []string, log bool) error {
	cmd := exec.Command(command, args...)
	if contextDir != "" {
		cmd.Dir = contextDir
	}
	acmd := aexec.New(cmd, 300)
	err := acmd.Run(true)
	if err != nil {
		return err
	}
	done := make(chan struct{})
	go func() {
		for b := range acmd.StdoutStream() {
			if log {
				fmt.Println(string(b))
			}
		}
		done <- struct{}{}
	}()
	go func() {
		for b := range acmd.StderrStream() {
			if log {
				fmt.Println(string(b))
			}
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

	if acmd.ExitCode() != 0 {
		return fmt.Errorf("exited with code %d", acmd.ExitCode())
	}

	return nil
}

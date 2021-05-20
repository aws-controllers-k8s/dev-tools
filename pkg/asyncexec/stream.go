package asyncexec

import (
	"fmt"
	"os"
	"os/exec"
)

// StreamCommand executes a given command in a context directory and streams
// the outputs to their according stdeout/stderr.
func StreamCommand(workDir string, command string, args []string) error {
	cmd := exec.Command(command, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	acmd := New(cmd, 8)
	err := acmd.Run()
	if err != nil {
		return err
	}

	// done is used to wait for stream readers to finish before exiting the function.
	done := make(chan struct{})

	go func() {
		for b := range acmd.StdoutStream() {
			_, err := os.Stdout.Write(b)
			if err != nil {
				msg := fmt.Sprintf("failed to write to Stdout: %v", err)
				// should never happen, just panic.
				panic(msg)
			}
		}
		done <- struct{}{}
	}()
	go func() {
		for b := range acmd.StderrStream() {
			_, err := os.Stderr.Write(b)
			if err != nil {
				msg := fmt.Sprintf("failed to write to Stderr: %v", err)
				// should never happen, just panic.
				panic(msg)
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

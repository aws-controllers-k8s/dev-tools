// +build !windows

package asyncexec_test

import (
	"fmt"
	"os/exec"

	"github.com/aws-controllers-k8s/dev-tools/pkg/asyncexec"
)

func ExampleCmd_Run_withNoStream() {
	cmd := asyncexec.New(exec.Command("echo", "Hello ACK"), 16)
	cmd.Run()
	cmd.Wait()
	fmt.Println(cmd.ExitCode())
	// Output: 0
}

func ExampleCmd_Run_withStream() {
	cmd := asyncexec.New(exec.Command("echo", "Hello ACK"), 16)
	cmd.Run()

	done := make(chan struct{})
	go func() {
		for b := range cmd.StdoutStream() {
			fmt.Println(string(b))

		}
		done <- struct{}{}
	}()
	go func() {
		for b := range cmd.StderrStream() {
			fmt.Println(string(b))

		}
		done <- struct{}{}
	}()

	defer func() { _, _ = <-done, <-done }()

	cmd.Wait()
	// Output: Hello ACK
}

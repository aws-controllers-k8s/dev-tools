package asyncexec

import (
	"bufio"
	"os/exec"
)

// New instantiate a new Cmd object.
func New(cmd *exec.Cmd, buff int) *Cmd {
	return &Cmd{
		cmd:      cmd,
		stopCh:   make(chan struct{}),
		stdoutCh: make(chan []byte, buff),
		stderrCh: make(chan []byte, buff),
	}
}

// Cmd is a wrapper arround exec.Cmd. Mainly used to execute
// command asynchronously with/or without output stream.
type Cmd struct {
	cmd *exec.Cmd

	stopCh   chan struct{}
	stdoutCh chan []byte
	stderrCh chan []byte
}

// Run runs the command. if streamOutput is true, it will spin
// two goroutine responsible of streaming the stdout and stderr
func (c *Cmd) Run() error {
	cmdStdoutReader, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdoutScanner := bufio.NewScanner(cmdStdoutReader)

	cmdStderrReader, err := c.cmd.StderrPipe()
	if err != nil {
		return err
	}
	stderrScanner := bufio.NewScanner(cmdStderrReader)

	// Goroutine for stdout
	go func() {
		defer close(c.stdoutCh)
		for stdoutScanner.Scan() {
			bytes := stdoutScanner.Bytes()
			c.stdoutCh <- bytes
		}
	}()

	// Goroutine for stderr
	go func() {
		defer close(c.stderrCh)
		for stderrScanner.Scan() {
			bytes := stderrScanner.Bytes()
			c.stderrCh <- bytes
		}
	}()

	err = c.cmd.Start()
	if err != nil {
		return err
	}

	// listening for stop signal
	go func() {
		<-c.stopCh
		c.cmd.Process.Kill()
	}()

	return nil
}

// Exited returns true if the command exited, false otherwise.
func (c *Cmd) Exited() bool {
	return c.cmd.ProcessState.Exited()
}

// ExitCode returns the command process exit code.
func (c *Cmd) ExitCode() int {
	return c.cmd.ProcessState.ExitCode()
}

// StdoutStream returns a channel streaming the command Stdout.
func (c *Cmd) StdoutStream() <-chan []byte {
	return c.stdoutCh
}

// StderrStream returns a channel streaming the command Stderr.
func (c *Cmd) StderrStream() <-chan []byte {
	return c.stderrCh
}

// Wait blocks until the command exits
func (c *Cmd) Wait() error {
	return c.cmd.Wait()
}

// Stop signals the Wrapper to kill the process running the command.
func (c *Cmd) Stop() {
	c.stopCh <- struct{}{}
}

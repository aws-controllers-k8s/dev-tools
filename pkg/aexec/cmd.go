package aexec

import (
	"bufio"
	"os/exec"
)

func New(cmd *exec.Cmd, buff int) *Cmd {
	return &Cmd{
		cmd:      cmd,
		stopChan: make(chan struct{}),
		StdoutCh: make(chan []byte, buff),
		StderrCh: make(chan []byte, buff),
	}
}

type Cmd struct {
	cmd *exec.Cmd

	stopChan chan struct{}

	Stdout [][]byte
	Stderr [][]byte

	StdoutCh chan []byte
	StderrCh chan []byte
}

func (c *Cmd) Run(streamOutput bool) error {
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
		defer close(c.StdoutCh)
		for stdoutScanner.Scan() {
			bytes := stdoutScanner.Bytes()
			c.Stdout = append(c.Stdout, bytes)
			if streamOutput {
				c.StdoutCh <- bytes
			}
		}
	}()

	// Goroutine for stderr
	go func() {
		defer close(c.StderrCh)
		for stderrScanner.Scan() {
			bytes := stderrScanner.Bytes()
			c.Stderr = append(c.Stderr, bytes)
			if streamOutput {
				c.StderrCh <- bytes
			}
		}
	}()

	err = c.cmd.Start()
	if err != nil {
		return err
	}

	// listening for stop signal
	go func() {
		<-c.stopChan
		c.cmd.Process.Kill()
	}()

	return nil
}

func (c *Cmd) Exited() bool {
	return c.cmd.ProcessState.Exited()
}

func (c *Cmd) ExitCode() int {
	return c.cmd.ProcessState.ExitCode()
}

func (c *Cmd) StdoutStream() <-chan []byte {
	return c.StdoutCh
}

func (c *Cmd) StderrStream() <-chan []byte {
	return c.StderrCh
}

func (c *Cmd) Wait() error {
	return c.cmd.Wait()
}

func (c *Cmd) Stop() {
	c.stopChan <- struct{}{}
}

type CmdResult struct {
	Service string
	Success bool
	Error   error
	Stdout  [][]byte
	Stderr  [][]byte
}

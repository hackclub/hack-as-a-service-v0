package main

import (
	"bufio"
	"context"
	"log"
	"os/exec"
)

var cmdExecId int = 0

/// A (possibly long running) command execution
type CmdExec struct {
	id         int
	cmd        *exec.Cmd
	stdoutChan chan string
	stderrChan chan string
	done       chan int
}

func (c *CmdExec) Stdout() chan string {
	return c.stdoutChan
}

func (c *CmdExec) Stderr() chan string {
	return c.stderrChan
}

func (c *CmdExec) Done() chan int {
	return c.done
}

func (c *CmdExec) Id() int {
	return c.id
}

func (c *CmdExec) Start() error {
	errChan := make(chan error)

	go func() {
		stdout, err := c.cmd.StdoutPipe()
		if err != nil {
			errChan <- err
			return
		}
		stderr, err := c.cmd.StderrPipe()
		if err != nil {
			errChan <- err
			return
		}
		err = c.cmd.Start()
		errChan <- err
		if err != nil {
			return
		}

		stdoutDone := make(chan bool)
		stderrDone := make(chan bool)

		go func() {
			defer func() {
				close(c.stdoutChan)
				stdoutDone <- true
			}()
			s := bufio.NewScanner(stdout)
			for s.Scan() {
				c.stdoutChan <- s.Text()
			}
		}()
		go func() {
			defer func() {
				close(c.stderrChan)
				stderrDone <- true
			}()
			s := bufio.NewScanner(stderr)
			for s.Scan() {
				c.stderrChan <- s.Text()
			}
		}()

		<-stdoutDone
		<-stderrDone
		err = c.cmd.Wait()
		if err == nil {
			c.done <- 0
			return
		}
		switch err := err.(type) {
		case *exec.ExitError:
			c.done <- err.ExitCode()
		default:
			log.Printf("Error: %+v\n", err)
			c.done <- 0
		}
	}()

	err := <-errChan
	log.Printf("Started command, err = %+v\n", err)

	if err != nil {
		log.Printf("Error[start]: %+v\n", err)
	}

	return err
}

func nextId() int {
	id := cmdExecId
	cmdExecId++
	return id
}

func NewCmdExec(ctx context.Context, name string, args []string) CmdExec {
	c := exec.CommandContext(ctx, name, args...)
	id := nextId()
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	doneChan := make(chan int)

	return CmdExec{
		id: id, cmd: c,
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
		done:       doneChan,
	}
}

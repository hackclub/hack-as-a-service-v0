package main

import (
	"bufio"
	"context"
	"log"
	"os/exec"

	"github.com/google/uuid"
)

/// A (possibly long running) command execution
type CmdExec struct {
	Id         uuid.UUID
	cmd        *exec.Cmd
	StdoutChan chan string
	StderrChan chan string
	Done       chan int
}

func (c *CmdExec) Start() error {
	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error[start]: %+v\n", err)
		return err
	}
	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		log.Printf("Error[start]: %+v\n", err)
		return err
	}
	err = c.cmd.Start()
	if err != nil {
		log.Printf("Error[start]: %+v\n", err)
		return err
	}

	go func() {
		go func() {
			defer close(c.StdoutChan)
			s := bufio.NewScanner(stdout)
			for s.Scan() {
				c.StdoutChan <- s.Text()
			}
		}()
		go func() {
			defer close(c.StderrChan)
			s := bufio.NewScanner(stderr)
			for s.Scan() {
				c.StderrChan <- s.Text()
			}
		}()

		err = c.cmd.Wait()
		if err == nil {
			c.Done <- 0
			return
		}
		switch err := err.(type) {
		case *exec.ExitError:
			c.Done <- err.ExitCode()
		default:
			log.Printf("Error: %+v\n", err)
			c.Done <- 0
		}
	}()

	return err
}

func NewCmdExec(ctx context.Context, name string, args []string) CmdExec {
	c := exec.CommandContext(ctx, name, args...)
	id, err := uuid.NewRandom()
	if err != nil {
		// should *not* happen
		log.Fatalln(err)
	}
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	doneChan := make(chan int)

	return CmdExec{
		Id: id, cmd: c,
		StdoutChan: stdoutChan,
		StderrChan: stderrChan,
		Done:       doneChan,
	}
}

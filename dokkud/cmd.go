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
	id          int
	cmd         *exec.Cmd
	stdoutChan  chan string
	stdoutStart chan bool
	stdoutDone  chan bool
	stderrChan  chan string
	stderrDone  chan bool
	stderrStart chan bool
	done        chan int
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
	err := c.cmd.Start()
	c.stdoutStart <- true
	c.stderrStart <- true
	log.Printf("Started command, err = %+v\n", err)

	// Done writer
	go (func() {
		<-c.stdoutDone
		<-c.stderrDone
		err := c.cmd.Wait()
		if err != nil {
			switch err := err.(type) {
			case *exec.ExitError:
				c.done <- err.ProcessState.ExitCode()
			default:
				c.done <- 0
				log.Printf("Error[status]: %+v\n", err)
			}
		}
		c.done <- 0
	})()

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
	stdoutDone := make(chan bool)
	stdoutStart := make(chan bool)
	stderrChan := make(chan string)
	stderrDone := make(chan bool)
	stderrStart := make(chan bool)
	doneChan := make(chan int)

	// Stdout writer
	go (func() {
		defer close(stdoutChan)
		defer func() { stdoutDone <- true }()
		stdout, err := c.StdoutPipe()
		if err != nil {
			log.Printf("Error[stdout1]: %+v\n", err)
			return
		}
		<-stdoutStart
		log.Println("Started stdout reader")
		r := bufio.NewScanner(stdout)
		for r.Scan() {
			line := r.Text()
			log.Printf("Line from stdout goroutine: %s\n", line)
			stdoutChan <- line
		}
		err = r.Err()
		if err != nil {
			log.Printf("Error[stdout2]: %+v\n", err)
		}
	})()

	// Stderr writer
	go (func() {
		defer close(stderrChan)
		defer func() { stderrDone <- true }()
		stderr, err := c.StderrPipe()
		if err != nil {
			log.Printf("Error[stderr1]: %+v\n", err)
			return
		}
		<-stderrStart
		log.Println("Started stderr reader")
		r := bufio.NewScanner(stderr)
		for r.Scan() {
			stderrChan <- r.Text()
		}
		if err != nil {
			log.Printf("Error[stderr2]: %+v\n", err)
		}
	})()

	return CmdExec{
		id: id, cmd: c,
		stdoutChan: stdoutChan, stdoutDone: stdoutDone, stdoutStart: stdoutStart,
		stderrChan: stderrChan, stderrDone: stderrDone, stderrStart: stderrStart,
		done: doneChan,
	}
}

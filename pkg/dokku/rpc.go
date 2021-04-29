package dokku

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"go.lsp.dev/jsonrpc2"
)

type DokkuConn struct {
	socketPath string
	conn       jsonrpc2.Conn
}

var commandOutputs map[uuid.UUID]map[StreamingCommandOutput]struct{} = make(map[uuid.UUID]map[StreamingCommandOutput]struct{})

func CreateOutput(execId uuid.UUID) (StreamingCommandOutput, error) {
	output := StreamingCommandOutput{
		StdoutChan: make(chan string),
		StderrChan: make(chan string),
		StatusChan: make(chan int),
	}
	err := AddCommandOutput(execId, output)
	if err != nil {
		return StreamingCommandOutput{}, err
	}
	return output, nil
}

func AddCommandOutput(execId uuid.UUID, output StreamingCommandOutput) error {
	if outputs, ok := commandOutputs[execId]; ok {
		outputs[output] = struct{}{}
		return nil
	}

	return fmt.Errorf("execution ID %d does not exist", execId)
}

func RemoveCommandOutput(execId uuid.UUID, output StreamingCommandOutput) {
	if outputs, ok := commandOutputs[execId]; ok {
		delete(outputs, output)
	}
}

// Exponentially backoff network errors
func retriedNetworkFunc(f func() (interface{}, error)) (interface{}, error) {
	// maximum backoff
	maxBackoff := 30
	currentBackoff := 1
	currentBackoffCounter := 0
	res, err := f()
	for {
		if err == nil {
			return res, err
		}
		switch err.(type) {
		case net.Error:
			// fallthrough to backoff
		default:
			return res, err
		}
		// backoff
		log.Printf("Network error, waiting %d seconds before trying again (attempt %d)\n", currentBackoff, currentBackoffCounter)
		time.Sleep(time.Duration(currentBackoff) * time.Second)
		currentBackoffCounter++
		if currentBackoffCounter == 2 {
			if currentBackoff == maxBackoff {
				break
			}
			currentBackoffCounter = 0
			currentBackoff *= 2
			if currentBackoff > maxBackoff {
				currentBackoff = maxBackoff
			}
		}
		res, err = f()
	}
	log.Printf("Network error, given up on retries (backoff time %d, attempt %d)\n", currentBackoff, currentBackoffCounter)
	return res, err
}

func (conn *DokkuConn) RunCommand(ctx context.Context, args []string) (string, error) {
	var res string
	_, err := conn.conn.Call(ctx, "command", args, &res)
	if err != nil {
		// fmt.Printf("Error: %+v\n", err)
		// FIXME: seems to reconnect every 2nd call
		if strings.Contains(err.Error(), "use of closed network connection") {
			err = conn.reconnect(ctx)
			if err != nil {
				return "", err
			}
			return conn.RunCommand(ctx, args)
		} else {
			return "", err
		}
	}

	return res, nil
}

type StreamingCommandOutput struct {
	StdoutChan chan string
	StderrChan chan string
	StatusChan chan int
}

type StreamingCommandResult struct {
	ExecId uuid.UUID
}

func (conn *DokkuConn) RunStreamingCommand(ctx context.Context, args []string) (StreamingCommandResult, error) {
	var res StreamingCommandResult
	_, err := conn.conn.Call(ctx, "streamingCommand", args, &res)
	if err != nil {
		// fmt.Printf("Error: %+v\n", err)
		// FIXME: seems to reconnect every 2nd call
		if strings.Contains(err.Error(), "use of closed network connection") {
			err = conn.reconnect(ctx)
			if err != nil {
				return StreamingCommandResult{}, err
			}
			return conn.RunStreamingCommand(ctx, args)
		} else {
			return StreamingCommandResult{}, err
		}
	}

	commandOutputs[res.ExecId] = make(map[StreamingCommandOutput]struct{})
	return res, nil
}

type lineMessage struct {
	ExecId uuid.UUID
	Line   string
}

type statusMessage struct {
	ExecId uuid.UUID
	Status int
}

func stdoutHandler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var msg lineMessage
	err := json.Unmarshal(req.Params(), &msg)
	if err != nil {
		return err
	}
	log.Printf("[stdout][%s] %s", msg.ExecId, msg.Line)
	for output := range commandOutputs[msg.ExecId] {
		output.StdoutChan <- msg.Line
	}
	// update db in background
	go func() {
		var build db.Build
		result := db.DB.First(&build, "exec_id = ?", msg.ExecId)
		// ignore errors, nothing we can do
		if result.Error != nil {
			return
		}
		build.Stdout += msg.Line + "\n"
		db.DB.Save(&build)
		log.Println("DB updated")
	}()
	return nil
}

func stderrHandler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var msg lineMessage
	err := json.Unmarshal(req.Params(), &msg)
	if err != nil {
		return err
	}
	log.Printf("[stderr][%s] %s", msg.ExecId, msg.Line)
	for output := range commandOutputs[msg.ExecId] {
		output.StderrChan <- msg.Line
	}
	// update db in background
	go func() {
		var build db.Build
		result := db.DB.First(&build, "exec_id = ?", msg.ExecId)
		// ignore errors, nothing we can do
		if result.Error != nil {
			return
		}
		build.Stderr += msg.Line + "\n"
		db.DB.Save(&build)
		log.Println("DB updated")
	}()
	return nil
}

func doneHandler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var msg statusMessage
	err := json.Unmarshal(req.Params(), &msg)
	if err != nil {
		return err
	}
	log.Printf("[%s] exited with status %d\n", msg.ExecId, msg.Status)
	for output := range commandOutputs[msg.ExecId] {
		output.StatusChan <- msg.Status
	}
	// update db in background
	go func() {
		var build db.Build
		result := db.DB.First(&build, "exec_id = ?", msg.ExecId)
		// ignore errors, nothing we can do
		if result.Error != nil {
			return
		}
		build.Status = msg.Status
		build.Running = false
		build.EndedAt = time.Now()
		db.DB.Save(&build)
		log.Println("DB updated")
	}()
	return nil
}

func notificationHandler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	switch req.Method() {
	case "commandStdout":
		return stdoutHandler(ctx, reply, req)
	case "commandStderr":
		return stderrHandler(ctx, reply, req)
	case "commandDone":
		return doneHandler(ctx, reply, req)
	default:
		return nil
	}
}

// Reconnects a dokku connection
func (conn *DokkuConn) reconnect(ctx context.Context) error {
	if conn.conn != nil {
		conn.conn.Close()
		<-conn.conn.Done()
	}
	stream, err := retriedNetworkFunc(func() (interface{}, error) { return net.Dial("unix", conn.socketPath) })
	if err != nil {
		return err
	}
	stream2 := jsonrpc2.NewStream(stream.(net.Conn))
	jconn := jsonrpc2.NewConn(stream2)
	jconn.Go(ctx, notificationHandler)
	conn.conn = jconn
	return nil
}

// Connects to the default dokku socket
func DokkuConnect(ctx context.Context) (*DokkuConn, error) {
	dconn := DokkuConn{socketPath: "/var/run/dokku-daemon/dokkud.sock"}
	err := dconn.reconnect(ctx)
	return &dconn, err
}

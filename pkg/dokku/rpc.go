package dokku

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku/util"
	"go.lsp.dev/jsonrpc2"
	"gorm.io/gorm"
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

type eventLine struct {
	Stream string
	Output string
}

type EventArgs struct {
	Event   string
	AppName string
}

var eventChannels map[chan EventArgs]struct{} = make(map[chan EventArgs]struct{})

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
		line := eventLine{Stream: "stdout", Output: msg.Line}
		out, err := json.Marshal(line)
		if err != nil {
			return
		}
		db.DB.Model(&build).Update("events", gorm.Expr("array_append(events, ?)", string(out)))
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
		line := eventLine{Stream: "stderr", Output: msg.Line}
		out, err := json.Marshal(line)
		if err != nil {
			return
		}
		db.DB.Model(&build).Update("events", gorm.Expr("array_append(events, ?)", string(out)))
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
		line := eventLine{Stream: "status", Output: strconv.Itoa(msg.Status)}
		out, err := json.Marshal(line)
		if err != nil {
			return
		}
		build.Events = append(build.Events, string(out))
		build.Status = msg.Status
		build.Running = false
		build.EndedAt = time.Now()
		db.DB.Model(&build).Updates(&build)
		log.Println("DB updated")
	}()
	return nil
}

func CreateEventChannel() chan EventArgs {
	ch := make(chan EventArgs)
	eventChannels[ch] = struct{}{}
	return ch
}

func DeleteEventChannel(ch chan EventArgs) {
	delete(eventChannels, ch)
}

func eventHandler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var args EventArgs
	err := json.Unmarshal(req.Params(), &args)
	if err != nil {
		return err
	}

	for ch := range eventChannels {
		ch <- args
	}

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
	case "event":
		return eventHandler(ctx, reply, req)
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
	stream, err := util.RetriedNetworkFunc(func() (interface{}, error) { return net.Dial("unix", conn.socketPath) })
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

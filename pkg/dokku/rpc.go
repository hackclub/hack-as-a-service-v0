package dokku

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"

	"go.lsp.dev/jsonrpc2"
)

type DokkuConn struct {
	socketPath string
	conn       jsonrpc2.Conn
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

type StreamingCommandResult struct {
	ExecId int
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

	return res, nil
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
	jconn.Go(ctx, func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
		var out map[string]interface{}
		json.Unmarshal(req.Params(), &out)
		log.Printf("Request method=%s params=%+v\n", req.Method(), out)
		return nil
	})
	conn.conn = jconn
	return nil
}

// Connects to the default dokku socket
func DokkuConnect(ctx context.Context) (*DokkuConn, error) {
	dconn := DokkuConn{socketPath: "/var/run/dokku-daemon/dokkud.sock"}
	err := dconn.reconnect(ctx)
	return &dconn, err
}

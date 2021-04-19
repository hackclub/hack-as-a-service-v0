package dokku

import (
	"context"
	"net"

	"go.lsp.dev/jsonrpc2"
)

type CommandParams struct {
	Command struct {
		Exe  string
		Args []string
	}
}

type CommandOutput struct {
	Stdout string
	Stderr string
}

type DokkuConn interface {
	// Run a command
	RunCommand(ctx context.Context, exe string, args []string) (CommandOutput, error)
}

type dokkuconn struct {
	conn jsonrpc2.Conn
}

func (conn *dokkuconn) RunCommand(ctx context.Context, exe string, args []string) (CommandOutput, error) {
	p := CommandParams{
		Command: struct {
			Exe  string
			Args []string
		}{
			exe, args,
		},
	}
	var res CommandOutput
	_, err := conn.conn.Call(ctx, "command", p, &res)
	if err != nil {
		return CommandOutput{}, err
	}
	return res, nil
}

// Connects to the default dokku socket
func DokkuConnect(ctx context.Context) (DokkuConn, error) {
	stream, err := net.Dial("unix", "/var/run/dokku-daemon/dokkud.sock")
	if err != nil {
		return nil, err
	}
	stream2 := jsonrpc2.NewStream(stream)
	conn := jsonrpc2.NewConn(stream2)
	conn.Go(ctx, jsonrpc2.MethodNotFoundHandler)
	return &dokkuconn{conn}, nil
}

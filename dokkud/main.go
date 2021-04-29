package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"os/exec"

	"go.lsp.dev/jsonrpc2"
)

func simpleCmdHandler(conn jsonrpc2.Conn, ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var args []string
	err := json.Unmarshal(req.Params(), &args)
	if err != nil {
		reply(ctx, nil, err)
		return nil
	}

	log.Printf("Request to execute command %s %s\n", "dokku", args)

	cmd := exec.Command("dokku", args...)

	stdout, err := cmd.Output()

	if err != nil {
		switch v := err.(type) {
		case *exec.ExitError:
			log.Printf("Error while running command: %s\n", v.Stderr)
			reply(ctx, nil, errors.New(string(v.Stderr)))
			return nil
		default:
			log.Println(err.Error())
			reply(ctx, "", nil)
			return nil
		}
	}

	output := string(stdout)

	err = reply(ctx, output, nil)
	if err != nil {
		return err
	}
	return nil
}

func streamingCmdHandler(conn jsonrpc2.Conn, ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var args []string
	err := json.Unmarshal(req.Params(), &args)
	if err != nil {
		reply(ctx, nil, err)
		return nil
	}

	log.Printf("Request to execute [streaming] command %s %s\n", "dokku", args)

	cmd := NewCmdExec(ctx, "dokku", args)

	err = cmd.Start()
	if err != nil {
		reply(ctx, nil, err)
		return nil
	}

	// Notifier
	go func() {
	loop:
		for {
			select {
			case line, ok := <-cmd.StdoutChan:
				if !ok {
					continue
				}
				conn.Notify(ctx, "commandStdout", map[string]interface{}{
					"execId": cmd.Id,
					"line":   line,
				})
			case line, ok := <-cmd.StderrChan:
				if !ok {
					continue
				}
				conn.Notify(ctx, "commandStderr", map[string]interface{}{
					"execId": cmd.Id,
					"line":   line,
				})
			case status := <-cmd.Done:
				conn.Notify(ctx, "commandDone", map[string]interface{}{
					"execId": cmd.Id,
					"status": status,
				})
				break loop
			}
		}
	}()

	err = reply(ctx, map[string]interface{}{
		"execId": cmd.Id,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

func mainHandler(conn jsonrpc2.Conn, ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	switch req.Method() {
	case "command":
		return simpleCmdHandler(conn, ctx, reply, req)
	case "streamingCommand":
		return streamingCmdHandler(conn, ctx, reply, req)
	default:
		return jsonrpc2.MethodNotFoundHandler(ctx, reply, req)
	}
}

func main() {
	path := "/var/run/dokku-daemon/dokkud.sock"
	os.Remove(path)

	log.Println("Listening on", path)
	ln, err := net.Listen("unix", path)
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()
	log.Println("Starting connection")
	for {
		c2, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		s := jsonrpc2.NewStream(c2)
		c := jsonrpc2.NewConn(s)
		c.Go(context.Background(), jsonrpc2.AsyncHandler(
			func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
				return mainHandler(c, ctx, reply, req)
			},
		))
	}
}

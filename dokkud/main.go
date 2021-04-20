package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"

	"go.lsp.dev/jsonrpc2"
)

func handler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	if req.Method() != "command" {
		reply(ctx, nil, errors.New("unsupported method"))
		return nil
	}

	var args []string
	err := json.Unmarshal(req.Params(), &args)
	if err != nil {
		reply(ctx, nil, err)
		return nil
	}

	log.Printf("Request to execute command %s %s", "dokku", args)

	cmd := exec.Command("dokku", args...)

	stdout, err := cmd.Output()

	if err != nil {
		switch v := err.(type) {
		case *exec.ExitError:
			log.Printf("Error while running command: %s", v.Stderr)
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

func main() {
	path := "/var/run/dokku-daemon/dokkud.sock"
	os.Remove(path)

	log.Println("Listening on", path)
	log.Fatalln(jsonrpc2.ListenAndServe(
		context.Background(), "unix", path, jsonrpc2.HandlerServer(jsonrpc2.AsyncHandler(handler)), 0,
	))
}

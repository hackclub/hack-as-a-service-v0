package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"go.lsp.dev/jsonrpc2"

	"github.com/hackclub/hack-as-a-service/dokku"
)

func handler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	if req.Method() != "command" {
		return errors.New("unsupported method")
	}

	var params dokku.CommandParams
	err := json.Unmarshal(req.Params(), &params)
	if err != nil {
		return err
	}

	log.Printf("Request to execute command %s %s", params.Command.Exe, strings.Join(params.Command.Args, " "))

	cmd := exec.Command(params.Command.Exe, params.Command.Args...)

	stderr := []byte{}
	stdout, err := cmd.Output()

	if err != nil {
		switch v := err.(type) {
		case *exec.ExitError:
			log.Printf("Error while running command: %s", v.Stderr)
			stderr = v.Stderr
		default:
			log.Println(err.Error())
			reply(ctx, dokku.CommandOutput{Stdout: "", Stderr: ""}, nil)
			return nil
		}
	}

	output := dokku.CommandOutput{
		Stdout: string(stdout),
		Stderr: string(stderr),
	}

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

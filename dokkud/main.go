package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"go.lsp.dev/jsonrpc2"

	"github.com/hackclub/hack-as-a-service/dokku"
)

func handler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	// log.Printf("Got connection! Method = %s\n", req.Method())

	if req.Method() != "command" {
		return errors.New("unsupported method")
	}

	var params dokku.CommandParams
	err := json.Unmarshal(req.Params(), &params)
	if err != nil {
		return err
	}

	log.Printf("Request to execute command %s %s", params.Command.Exe, strings.Join(params.Command.Args, " "))

	// log.Printf("Exe = %s, Args = %+v\n", params.Command.Exe, params.Command.Args)

	cmd := exec.Command(params.Command.Exe, params.Command.Args...)
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	stdout, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return err
	}
	stderr, err := ioutil.ReadAll(stderrPipe)
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	output := dokku.CommandOutput{
		Stdout: string(stdout),
		Stderr: string(stderr),
	}

	// log.Printf("Output: %+v\n", output)

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

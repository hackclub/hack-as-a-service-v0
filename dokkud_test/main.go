package main

import (
	"context"
	"log"
	"net"

	"go.lsp.dev/jsonrpc2"
)

type Params struct {
	Command struct {
		Exe  string
		Args []string
	}
}

type Output struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func main() {
	conn, err := net.Dial("unix", "/var/run/dokkud.sock")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to socket")
	stream := jsonrpc2.NewStream(conn)
	conn2 := jsonrpc2.NewConn(stream)
	defer conn2.Close()
	log.Println("Setup jsonrpc2 conn")
	res := Output{}
	conn2.Go(context.Background(), jsonrpc2.MethodNotFoundHandler)
	_, err = conn2.Call(context.Background(), "command", Params{
		Command: struct {
			Exe  string
			Args []string
		}{
			Exe:  "dokku",
			Args: []string{"help"},
		},
	}, &res)
	log.Println("call done")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Result of dokku help:(stdout, stderr)\n%s\n%s\n", res.Stdout, res.Stderr)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/dokku/dokku/plugins/common"
	"github.com/hackclub/hack-as-a-service/pkg/dokku/util"
	"go.lsp.dev/jsonrpc2"
)

func main() {
	parts := strings.Split(os.Args[0], "/")
	trigger := parts[len(parts)-1]
	flag.Parse()

	var err error
	switch trigger {
	case "post-deploy":
		app := os.Args[1]
		err = SendEvent("post-deploy", app)
	default:
		err = fmt.Errorf("invalid plugin trigger call: %s", trigger)
	}

	if err != nil {
		common.LogFailWithError(err)
	}
}

type EventArgs struct {
	Event   string
	AppName string
}

func SendEvent(event string, appName string) error {
	ctx := context.Background()

	stream, err := util.RetriedNetworkFunc(func() (interface{}, error) { return net.Dial("unix", "/var/run/dokku-daemon/dokkud.sock") })
	if err != nil {
		return err
	}
	stream2 := jsonrpc2.NewStream(stream.(net.Conn))
	jconn := jsonrpc2.NewConn(stream2)
	jconn.Go(ctx, jsonrpc2.MethodNotFoundHandler)

	err = jconn.Notify(ctx, "event", EventArgs{Event: event, AppName: appName})
	if err != nil {
		return err
	}

	return nil
}

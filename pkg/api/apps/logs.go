package apps

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handleGETLogs(c *gin.Context) {
	dokku_conn := c.MustGet("dokkuconn").(dokku.DokkuConn)

	app_name := c.Param("id")

	// Get the app's container ID
	cid, err := dokku_conn.RunCommand(context.Background(), []string{"haas:cid", app_name})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	cid = strings.TrimSpace(cid)

	// Initialize a Docker API client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "Error creating Docker client"})
		return
	}

	// Open a log stream via the Docker API
	log_stream, err := cli.ContainerLogs(context.Background(), cid, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "10",
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Make a channel to hold the log stream
	log_chan := make(chan []byte)

	// A channel to close the client stream when the Docker API disconnects
	log_stream_closed := make(chan bool)

	// Listen for new logs in a goroutine
	go func() {
		for {
			logs := make([]byte, 100)
			_, err := log_stream.Read(logs)
			if err != nil {
				log_stream_closed <- true
				return
			}

			log_chan <- logs
		}
	}()

	// Close the log stream when done
	defer log_stream.Close()

	clientGone := c.Writer.CloseNotify()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case <-log_stream_closed:
			return false
		case logs := <-log_chan:
			w.Write(logs)
		}
		return true
	})
}

package apps

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handleGETLogs(c *gin.Context) {
	dokku_conn := c.MustGet("dokkuconn").(dokku.DokkuConn)
	_ = c.MustGet("user").(db.User)

	app_id := c.Param("id")

	// var app db.App
	// result := db.DB.Preload("Team").First(&app, "id = ?", app_id)
	// if result.Error != nil {
	// 	c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
	// 	return
	// }

	// c.String(200, app.Team.Name)
	// return

	// Get the app's container ID
	cid, err := dokku_conn.RunCommand(context.Background(), []string{"haas:cid", app_id})
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

	// A channel to close the client stream when the Docker API disconnects
	log_stream_closed := make(chan bool)

	stdout_reader, stdout_writer := io.Pipe()
	stderr_reader, stderr_writer := io.Pipe()

	// Demultiplex stream in a goroutine
	go func() {
		defer stdout_reader.Close()
		defer stderr_reader.Close()

		defer func() { log_stream_closed <- true }()

		stdcopy.StdCopy(stdout_writer, stderr_writer, log_stream)
	}()

	// Make a channel to hold the log stream
	log_chan := make(chan []byte)

	// Listen for new logs in a goroutine
	go func() {
		for {
			logs := make([]byte, 100)
			_, err := stdout_reader.Read(logs)
			if err != nil {
				return
			}

			log_chan <- logs
		}
	}()

	go func() {
		for {
			logs := make([]byte, 100)
			_, err := stderr_reader.Read(logs)
			if err != nil {
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

package apps

import (
	"context"
	"fmt"
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
	app_id, err := dokku_conn.RunCommand(context.Background(), []string{"haas:cid", app_name})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	app_id = strings.TrimSpace(app_id)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "Error creating Docker client"})
		return
	}

	log_stream, err := cli.ContainerLogs(context.Background(), app_id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "10",
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	defer log_stream.Close()

	c.Stream(func(w io.Writer) bool {
		logs := make([]byte, 100)
		log_stream.Read(logs)
		w.Write(logs)

		fmt.Println("loop")

		return true
	})

	fmt.Println("woot")
}

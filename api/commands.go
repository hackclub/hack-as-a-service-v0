package api

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/dokku"
)

func handleApiCommand(c *gin.Context) {
	conn := c.MustGet("dokkuconn").(dokku.DokkuConn)
	args := strings.Split(c.Query("command"), " ")

	res, err := conn.RunCommand(context.Background(), args)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "error": err})
	} else {
		c.JSON(200, gin.H{"status": "ok", "output": res})
	}
}

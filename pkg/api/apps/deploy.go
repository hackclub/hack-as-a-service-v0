package apps

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handlePOSTDeploy(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	var json struct {
		GitRepository string
	}

	err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	user := c.MustGet("user").(db.User)

	var app db.App
	result := db.DB.Joins("INNER JOIN team_users ON apps.team_id = team_users.team_id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	// Delete previous clone
	_, err = conn.RunCommand(c.Request.Context(), []string{
		"git:unlock", app.ShortName, "--force",
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	cmd, err := conn.RunStreamingCommand(c.Request.Context(), []string{
		"git:sync", "--build", app.ShortName, json.GitRepository,
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Cannot provide WebSocket since this is a POST request
	c.JSON(200, gin.H{"status": "ok", "execId": cmd.ExecId})
}

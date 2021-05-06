package apps

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

type Output struct {
	LetsEncryptEnabled bool
}

func handleGETLetsEncrypt(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	var app db.App
	result := db.DB.Joins("INNER JOIN team_users ON team_users.team_id = apps.team_id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	conn := c.MustGet("dokku_conn").(*dokku.DokkuConn)

	res, err := conn.RunCommand(c.Request.Context(), []string{"haas:letsencrypt-enabled", app.ShortName})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}
	enabled := strings.TrimSpace(res) == "true"

	output := Output{
		LetsEncryptEnabled: enabled,
	}
	c.JSON(200, gin.H{"status": "ok", "letsencrypt": output})
}

func handlePOSTLetsEncryptEnable(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	var app db.App
	result := db.DB.Joins("INNER JOIN team_users ON team_users.team_id = apps.team_id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	conn := c.MustGet("dokku_conn").(*dokku.DokkuConn)

	_, err = conn.RunCommand(c.Request.Context(), []string{"letsencrypt:enable", app.ShortName})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

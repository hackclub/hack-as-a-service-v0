package apps

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handlePOSTRestart(c *gin.Context) {
	user := c.MustGet("user").(db.User)
	app_id := c.Param("id")
	dokku_conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)

	var app db.App
	result := db.DB.Joins("JOIN team_users ON team_users.team_id = apps.team_id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", app_id, user.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	_, err := dokku_conn.RunCommand(c.Request.Context(), []string{"ps:restart", app.ShortName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

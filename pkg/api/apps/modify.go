package apps

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handlePATCHApp(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	var json struct {
		Name string
		// TODO: Rename the dokku app
	}
	err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	var app db.App
	result := db.DB.Joins("INNER JOIN team_users ON team_users.team_id = apps.id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	app.Name = json.Name

	result = db.DB.Save(&app)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "app": app})
}

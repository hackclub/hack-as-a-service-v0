package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETAuthedApps(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	var apps []db.App
	result := db.DB.
		Joins("INNER JOIN team_users ON team_users.team_id = apps.team_id").
		Where("team_users.user_id = ?", user.ID).Find(&apps)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "apps": apps})
	}
}

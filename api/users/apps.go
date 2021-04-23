package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
)

func handleGETUserApps(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	var apps []db.App
	result := db.DB.
		Joins("INNER JOIN team_users ON team_users.team_id = apps.team_id").
		Where("team_users.user_id = ?", uint(id)).Find(&apps)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "apps": apps})
	}
}

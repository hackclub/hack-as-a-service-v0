package apps

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETBuilds(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	limit := 50
	if qlimit := c.Query("limit"); qlimit != "" {
		limit2, err := strconv.Atoi(qlimit)
		if err == nil {
			limit = limit2
		}
	}

	var builds []db.Build
	result := db.DB.Order("started_at DESC").Limit(limit).
		Joins("INNER JOIN team_users ON team_users.team_id = apps.id").
		Joins("INNER JOIN apps ON apps.id = builds.app_id").
		Find(&builds, "builds.app_id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "builds": builds})
}

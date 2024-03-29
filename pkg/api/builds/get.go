package builds

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETBuild(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid build ID"})
		return
	}

	var build db.Build
	result := db.DB.
		Joins("INNER JOIN apps ON apps.id = builds.app_id").
		Joins("INNER JOIN team_users ON team_users.team_id = apps.team_id").
		First(&build, "builds.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid build ID"})
	} else {
		c.JSON(200, gin.H{"status": "ok", "build": build})
	}
}

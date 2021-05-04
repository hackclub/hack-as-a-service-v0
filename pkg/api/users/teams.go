package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETAuthedTeams(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	var teams []db.Team

	err := db.DB.Model(&user).Association("Teams").Find(&teams)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "teams": teams})
}

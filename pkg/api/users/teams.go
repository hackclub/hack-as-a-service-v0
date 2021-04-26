package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"gorm.io/gorm"
)

func handleGETUserTeams(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	var teams []db.Team

	err = db.DB.Model(&db.User{
		Model: gorm.Model{
			ID: uint(id),
		},
	}).Association("Teams").Find(&teams)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "teams": teams})
}

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

package apps

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
)

func handlePOSTApp(c *gin.Context) {
	var json struct {
		Name   string
		TeamID uint
	}

	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	// Check that the given team ID exists
	result := db.DB.First(&db.Team{}, "id = ?", json.TeamID)
	if result.Error != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid team ID"})
		return
	}

	// create in db
	app := db.App{Name: json.Name, TeamID: json.TeamID}
	result = db.DB.Create(&app)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "app": app})
	}
}

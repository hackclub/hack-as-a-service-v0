package apps

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
	"gorm.io/gorm"
)

func handlePUTAppUsers(c *gin.Context) {
	_db := c.MustGet("db").(*gorm.DB)
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	var json struct {
		Users []uint
	}
	err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	var app db.App
	result := _db.First(&app, "id = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
		return
	}
	app.Users = nil
	for _, user := range json.Users {
		app.Users = append(app.Users, db.User{Model: gorm.Model{ID: user}})
	}

	result = _db.Save(&app)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok"})
	}
}

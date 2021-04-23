package apps

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
	"gorm.io/gorm"
)

func handlePOSTApp(c *gin.Context) {
	var json struct {
		Name             string
		BillingAccountID uint
	}

	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	// create in db
	app := db.App{Name: json.Name, BillingAccount: db.BillingAccount{Model: gorm.Model{ID: json.BillingAccountID}}}
	result := db.DB.Create(&app)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "appID": app.ID})
	}
}

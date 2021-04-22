package db

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupAppsRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTApp)
	r.GET("/:id", handleGETApp)
	r.DELETE("/:id", handleDELETEApp)
	r.PUT("/:id/users", handlePUTAppUsers)
}

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
	app := App{Name: json.Name, BillingAccount: BillingAccount{Model: gorm.Model{ID: json.BillingAccountID}}}
	result := db.Create(&app)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "appID": app.ID})
	}
}

func handleGETApp(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	var app App
	result := db.First(&app, "id = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "app": app})
	}
}

func handleDELETEApp(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	result := db.Delete(&App{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok"})
	}
}

func handlePUTAppUsers(c *gin.Context) {
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

	var app App
	result := db.First(&app, "id = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
		return
	}
	app.Users = nil
	for _, user := range json.Users {
		app.Users = append(app.Users, User{Model: gorm.Model{ID: user}})
	}

	result = db.Save(&app)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok"})
	}
}

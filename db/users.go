package db

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func setupUserRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTUser)
	r.GET("/:id", handleGETUser)
	r.DELETE("/:id", handleDELETEUser)
	r.GET("/:id/apps", handleGETUserApps)
}

func handleDELETEUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	result := db.Delete(&User{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok"})
	}
}

func handleGETUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	var user User
	result := db.First(&user, "id = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "user": user})
	}
}

func handlePOSTUser(c *gin.Context) {
	var json struct {
		SlackUserID string
	}

	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	// create in db
	user := User{SlackUserID: json.SlackUserID}
	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "userID": user.ID})
	}
}

func handleGETUserApps(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	var apps []App
	result := db.Raw("SELECT apps.* FROM apps INNER JOIN user_apps ON user_apps.app_id = apps.id WHERE user_apps.user_id = ?", uint(id)).Scan(&apps)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "apps": apps})
	}
}

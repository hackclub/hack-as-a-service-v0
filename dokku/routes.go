package dokku

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var db *gorm.DB

func get_gsn() string {
	db_password := os.Getenv("POSTGRES_PASSWORD")
	return fmt.Sprintf("host=db user=postgres password=%s dbname=haas port=5432 sslmode=disable", db_password)
}

func SetupRoutes(r *gin.RouterGroup) error {
	_db, err := gorm.Open(postgres.Open(get_gsn()), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	err = db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	r.POST("/users", handlePOSTUser)
	r.GET("/users/:id", handleGETUser)

	return nil
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

package apps

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func IsValidAppName(appName string) error {
	if appName == "" {
		return errors.New("please specify an app to run the command on")
	}

	r, _ := regexp.Compile(`^[a-z0-9][^/:_A-Z\s]*$`)
	if r.MatchString(appName) {
		return nil
	}

	return errors.New("app name must begin with lowercase alphanumeric character, and cannot include uppercase characters, colons, underscores, or whitespace")
}

func handlePOSTApp(c *gin.Context) {
	dokku_conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	user := c.MustGet("user").(db.User)

	var json struct {
		Name      string
		ShortName string
		TeamID    uint
	}

	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	err = IsValidAppName(json.ShortName)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Check that the given team ID exists
	result := db.DB.Joins("JOIN team_users ON team_users.team_id = teams.id").
		First(&db.Team{}, "teams.id = ? AND team_users.user_id = ?", json.TeamID, user.ID)
	if result.Error != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid team ID"})
		return
	}

	// Create Dokku app
	_, err = dokku_conn.RunCommand(c.Request.Context(), []string{"apps:create", json.ShortName})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "Error provisioning app"})
		return
	}

	// create in db
	app := db.App{Name: json.Name, TeamID: json.TeamID, ShortName: json.ShortName}
	result = db.DB.Create(&app)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "app": app})
}

package teams

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api/util"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"gorm.io/gorm"
)

func handlePATCHTeam(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	user := c.MustGet("user").(db.User)

	var json struct {
		AddUsers    []uint
		RemoveUsers []uint
	}
	err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	json.AddUsers, json.RemoveUsers = util.RemoveDuplicates(json.AddUsers, json.RemoveUsers)

	var addUsers []db.User
	var removeUsers []db.User

	for _, u := range json.AddUsers {
		if u == user.ID {
			c.JSON(400, gin.H{"status": "error", "message": "You can't invite yourself"})
			return
		}

		addUsers = append(addUsers, db.User{
			Model: gorm.Model{
				ID: u,
			},
		})
	}

	for _, u := range json.RemoveUsers {
		removeUsers = append(removeUsers, db.User{
			Model: gorm.Model{
				ID: u,
			},
		})
	}

	var team db.Team
	result := db.DB.Joins("JOIN team_users ON teams.id = team_users.team_id").
		First(&team, "teams.id = ? AND team_users.user_id = ? AND NOT teams.personal", id, user.ID)
	if result.Error != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid team ID"})
		return
	}

	err = db.DB.Model(&team).Association("Users").Append(addUsers)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	err = db.DB.Model(&team).Association("Users").Delete(removeUsers)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

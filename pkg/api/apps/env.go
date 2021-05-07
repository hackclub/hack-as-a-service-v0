package apps

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api/util"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handleGETEnv(c *gin.Context) {
	user := c.MustGet("user").(db.User)
	dokku_conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	app_id := c.Param("id")

	var app db.App
	// result := db.DB.First(&app, "id = ?", app_id)
	result := db.DB.Joins("JOIN team_users ON team_users.team_id = apps.team_id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", app_id, user.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	env_string, err := dokku_conn.RunCommand(c.Request.Context(), []string{"config:export", "--format", "json", app.ShortName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	raw_env := make(map[string]string)
	err = json.Unmarshal([]byte(env_string), &raw_env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	env := make(map[string]string)

	for key, value := range raw_env {
		if !util.IsCoreEnvVariable(key) {
			env[key] = value
		}
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"env":    env,
	})
}

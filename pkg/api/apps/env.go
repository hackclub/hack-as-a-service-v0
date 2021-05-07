package apps

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api/util"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func getAppEnv(ctx context.Context, dokku *dokku.DokkuConn, app string, includeCore bool) (map[string]string, error) {
	env_string, err := dokku.RunCommand(ctx, []string{"config:export", "--format", "json", app})
	if err != nil {
		return nil, err
	}

	raw_env := make(map[string]string)
	err = json.Unmarshal([]byte(env_string), &raw_env)
	if err != nil {
		return nil, err
	}

	if includeCore {

		return raw_env, nil
	}

	env := make(map[string]string)

	for key, value := range raw_env {
		if !util.IsCoreEnvVariable(key) {
			env[key] = value
		}
	}

	return env, nil
}

func handleGETEnv(c *gin.Context) {
	user := c.MustGet("user").(db.User)
	dokku_conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	app_id := c.Param("id")

	var app db.App
	result := db.DB.Joins("JOIN team_users ON team_users.team_id = apps.team_id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", app_id, user.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	env, err := getAppEnv(c.Request.Context(), dokku_conn, app.ShortName, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"env":    env,
	})
}

func handlePUTEnv(c *gin.Context) {
	user := c.MustGet("user").(db.User)
	app_id := c.Param("id")
	dokku_conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)

	var json struct {
		Env map[string]string
	}

	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid JSON",
		})
		return
	}

	var app db.App
	result := db.DB.Joins("JOIN team_users ON team_users.team_id = apps.team_id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", app_id, user.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	env, err := getAppEnv(c.Request.Context(), dokku_conn, app.ShortName, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, env)
}

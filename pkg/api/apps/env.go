package apps

import (
	"context"
	"encoding/json"
	"fmt"
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

	// Make sure no core variables are being set
	for key := range json.Env {
		if !util.VerifyEnv(key) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Invalid key: %s", key),
			})
			return
		} else if util.IsCoreEnvVariable(key) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Can't set reserved environment variable: %s", key),
			})
			return
		}
	}

	// Fetch the current app environment
	env, err := getAppEnv(c.Request.Context(), dokku_conn, app.ShortName, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	toDelete := []string{}

	// Find environment variables to delete
	for key := range env {
		if _, ok := json.Env[key]; !ok {
			toDelete = append(toDelete, key)
		}
	}

	if len(toDelete) > 0 {
		_, err := dokku_conn.RunCommand(c.Request.Context(), append([]string{"config:unset", "--no-restart", app.ShortName}, toDelete...))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	args := append([]string{"config:set", "--no-restart", "--encoded", app.ShortName}, util.FormatEnv(json.Env)...)
	_, err = dokku_conn.RunCommand(c.Request.Context(), args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "success",
	})
}

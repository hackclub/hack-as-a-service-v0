package users

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/me", handleGETAuthed)
	r.GET("/me/apps", handleGETAuthedApps)
	r.GET("/me/teams", handleGETAuthedTeams)
}

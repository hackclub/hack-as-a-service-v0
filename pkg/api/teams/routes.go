package teams

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTTeam)
	r.GET("/:id", handleGETTeam)
	r.GET("/me", handleGETPersonalTeam)
	r.PATCH("/:id", handlePATCHTeam)
	r.GET("/:id/apps", handleGETTeamApps)
}

package teams

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTTeam)
	r.GET("/me", handleGETPersonalTeam)
	r.GET("/:id", handleGETTeam)
	r.PATCH("/:id", handlePATCHTeam)
	r.GET("/:id/expenses", handleGETExpenses)
	r.GET("/:id/apps", handleGETTeamApps)
}

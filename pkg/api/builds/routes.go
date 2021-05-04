package builds

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/:id", handleGETBuild)
	r.GET("/:id/logs", handleGETLogs)
}

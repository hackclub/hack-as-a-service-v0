package apps

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTApp)
	r.GET("/:id", handleGETApp)
	r.GET("/:id/logs", handleGETLogs)
	r.GET("/:id/build-logs", handleGETBuildLogs)
	r.POST("/:id/deploy", handlePOSTDeploy)
	r.DELETE("/:id", handleDELETEApp)
}

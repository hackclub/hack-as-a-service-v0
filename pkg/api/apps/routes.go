package apps

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTApp)
	r.GET("/:id", handleGETApp)
	r.POST("/:id/rename", handlePOSTRename)
	r.GET("/:id/builds", handleGETBuilds)
	r.GET("/:id/logs", handleGETLogs)
	r.POST("/:id/deploy", handlePOSTDeploy)
	r.DELETE("/:id", handleDELETEApp)
}

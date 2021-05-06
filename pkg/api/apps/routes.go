package apps

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTApp)
	r.GET("/:id", handleGETApp)
	r.PATCH("/:id", handlePATCHApp)
	r.GET("/:id/builds", handleGETBuilds)
	r.GET("/:id/logs", handleGETLogs)
	r.GET("/:id/stats", handleGETStats)
	r.GET("/:id/letsencrypt", handleGETLetsEncrypt)
	r.POST("/:id/letsencrypt/enable", handlePOSTLetsEncryptEnable)
	r.POST("/:id/deploy", handlePOSTDeploy)
	r.DELETE("/:id", handleDELETEApp)
}

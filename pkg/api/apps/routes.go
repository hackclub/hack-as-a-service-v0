package apps

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTApp)
	r.GET("/:id", handleGETApp)
	r.DELETE("/:id", handleDELETEApp)
}

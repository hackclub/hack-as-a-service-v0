package users

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTUser)
	r.GET("/", handleGETAuthed)
	r.GET("/:id", handleGETUser)
	r.DELETE("/:id", handleDELETEUser)
	r.GET("/:id/apps", handleGETUserApps)
}

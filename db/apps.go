package db

import (
	"github.com/gin-gonic/gin"
)

func setupAppsRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTApp)
	r.GET("/:id", handleGETApp)
	r.DELETE("/:id", handleDELETEApp)
}

func handlePOSTApp(c *gin.Context) {
	// TODO
}

func handleGETApp(c *gin.Context) {
	// TODO
}

func handleDELETEApp(c *gin.Context) {
	// TODO
}

package gh

import (
	"github.com/gin-gonic/gin"
)

// Setup routes associated with the github bot
func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/webhook", handleWebhook)
	r.POST("/install", handleInstall)
}

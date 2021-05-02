package gh

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/install", handleInstall)
	r.POST("/webhook", handleWebhook)
}

func handleInstall(c *gin.Context) {
	fmt.Println("Handled install!")
}

func handleWebhook(c *gin.Context) {
	fmt.Println("Handled webhook!")
}

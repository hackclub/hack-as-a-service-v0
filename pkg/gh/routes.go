package gh

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/install", handleInstall)
	r.POST("/webhook", handleWebhook)
}

func handleInstall(c *gin.Context) {
	fmt.Println("Started install!")
	installID := c.Request.URL.Query().Get("installation_id")
	genClient(installID)
	fmt.Println("Made client!")
}

func handleWebhook(c *gin.Context) {
	fmt.Println("Handled webhook!")
}

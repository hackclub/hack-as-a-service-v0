package billing

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTBillingAccount)
	r.GET("/:id", handleGETBillingAccount)
	r.DELETE("/:id", handleDELETEBillingAccount)
}

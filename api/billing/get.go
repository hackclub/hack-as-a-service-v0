package billing

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
	"gorm.io/gorm"
)

func handleGETBillingAccount(c *gin.Context) {
	_db := c.MustGet("db").(*gorm.DB)
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid billing account ID"})
		return
	}

	var account db.BillingAccount
	result := _db.First(&account, "id = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "billingAccount": account})
	}
}

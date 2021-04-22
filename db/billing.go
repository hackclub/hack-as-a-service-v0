package db

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func setupBillingAccountsRoutes(r *gin.RouterGroup) {
	r.POST("/", handlePOSTBillingAccount)
	r.GET("/:id", handleGETBillingAccount)
	r.DELETE("/:id", handleDELETEBillingAccount)
}

func handlePOSTBillingAccount(c *gin.Context) {
	var json struct {
		HNUserID string
	}

	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	// create in db
	account := BillingAccount{HNUserID: json.HNUserID}
	result := db.Create(&account)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "billingAccountID": account.ID})
	}
}

func handleGETBillingAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid billing account ID"})
		return
	}

	var account BillingAccount
	result := db.First(&account, "id = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "billingAccount": account})
	}
}

func handleDELETEBillingAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid billing account ID"})
		return
	}

	result := db.Delete(&BillingAccount{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok"})
	}
}

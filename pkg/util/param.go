package util

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Get a query a parameter from a POST request that is required.
// If there is no value than an error is returned
func ReqGetQuery(name string, c *gin.Context) (string, error) {
	val, exists := c.GetQuery(name)
	if !exists {
		return "", errors.New("Failed to get value for required query parameter " + name)
	}
	return val, nil
}

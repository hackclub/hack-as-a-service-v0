package api

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func getApiKey() string {
	if key, ok := os.LookupEnv("API_KEY"); ok {
		return key
	} else {
		return "testinghaasapikey"
	}
}

func RequireBearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.Query("api_key")

		if api_key == "" {
			// Get from auth header if possible
			if auth_header := c.GetHeader("Authorization"); auth_header != "" {
				if strings.HasPrefix(auth_header, "Bearer ") {
					api_key = strings.TrimPrefix(auth_header, "Bearer ")
				}
			}
		}

		if api_key != getApiKey() {
			c.AbortWithStatusJSON(401, gin.H{"status": "error", "message": "Invalid API key"})
		} else {
			c.Next()
		}
	}
}

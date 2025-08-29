package handlers

import (
	"inventory_management/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	adminUser := getenvDefault("ADMIN_USERNAME", "admin")
	adminPass := getenvDefault("ADMIN_PASSWORD", "password")
	if req.Username == adminUser && req.Password == adminPass {
		token, err := middleware.GenerateJWT(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		v = def
	}
	return v
}

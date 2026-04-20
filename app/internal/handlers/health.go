package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func Chaos(c *gin.Context) {
	os.Exit(1)
}
package handlers

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func HealthCheck(c *gin.Context) {
	if RedisClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "no redis client"})
		return
	}

	if err := RedisClient.Ping(c).Err(); err != nil {
    c.JSON(http.StatusOK, gin.H{
        "status": "degraded",
        "redis":  "down",
    })
			return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Chaos(c *gin.Context) {
	os.Exit(1)
}
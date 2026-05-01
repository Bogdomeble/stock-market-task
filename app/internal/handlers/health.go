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
	// // we send request with 200 code to let nginx know that operation succeded

	// c.JSON(http.StatusOK, gin.H{"message": "chaos initiated"})
 //    go func() {
 //        time.Sleep(100 * time.Millisecond) // maybe not needed, slows down the requests per minute metric
 //        os.Exit(1)
 //    }()

 // we use nginx config to fix this problem

	os.Exit(1)
}

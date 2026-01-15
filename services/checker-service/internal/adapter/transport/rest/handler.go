package rest

import "github.com/gin-gonic/gin"

func InitHealthCheck() *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "worker is running"})
	})
	return r
}

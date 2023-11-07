package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "REST API using gin and mongodb",
		})
	})

	configs.ConnectDB()

	routes.BookRoute(router)

	router.Run("localhost:8080")
}

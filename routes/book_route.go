package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoute(router *gin.Engine) {
	router.POST("/books", controllers.CreateBook())
	router.GET("/books/:id", controllers.GetBook())
	router.PUT("/books/:id", controllers.EditBook())
	router.DELETE("/books/:id", controllers.DeleteBook())
	router.GET("/books", controllers.GetAllBooks())
}

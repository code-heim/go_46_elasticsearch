package main

import (
	"gin_elasticsearch/controllers"
	"gin_elasticsearch/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())

	r.LoadHTMLGlob("templates/**/*")

	models.ConnectDatabase()
	models.DBMigrate()

	r.GET("/blogs", controllers.BlogsIndex)
	r.GET("/blogs/:id", controllers.BlogsShow)

	log.Println("Server started!")
	r.Run() // Default Port 8080
}

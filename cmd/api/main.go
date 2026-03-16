package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"order-management-system/internal/database"
	"order-management-system/internal/handlers"
)

func main() {

	err := database.ConnectMongo()

	if err != nil {
		log.Fatal("Mongo Connection Failed", err)
	}
	router := gin.Default()

	router.POST("/orders", handlers.CreateOrder)

	router.Run(":8080")
}

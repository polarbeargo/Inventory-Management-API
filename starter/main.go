package main

import (
	"log"

	"inventory_management/database"
	"inventory_management/routes"
)

func main() {

	database.InitRedis()
	database.InitDatabase()
	log.Println("Server successfully connected to the database and seeded data.")
	router := routes.SetupRoutes()
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

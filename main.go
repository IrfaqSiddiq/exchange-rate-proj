package main

import (
	"log"
	"os"
	"project_first/routes"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file -> ", err)
	}
	//setup routes
	r := routes.SetupRouter()
	// running
	r.Run(":" + os.Getenv("PORT"))
}

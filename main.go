package main

import (
	"deck/config"
	"deck/database"
	"deck/routes"
	"os"
)

func main() {
	config.LoadEnv()

	database.InitDB()

	r := routes.SetupRoutes()

	r.Run(":" + os.Getenv("APP_PORT"))
}

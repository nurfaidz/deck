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

	r.Run("0.0.0.0:" + os.Getenv("APP_PORT"))
}

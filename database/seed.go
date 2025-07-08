package database

import (
	"deck/helpers"
	"deck/models"
	"fmt"
)

func SeedUser() {

	if DB == nil {
		fmt.Println("Error: DB is nil in SeedUser()")
		return
	}

	var user models.User

	result := DB.First(&user, "username = ?", "admin")
	if result.RowsAffected > 0 {
		fmt.Println("User already exists")

		return
	}

	hashedPassword := helpers.HashPassword("password")

	user = models.User{
		Username: "admin",
		Email:    "admin@gmail.com",
		Password: hashedPassword,
	}

	if err := DB.Create(&user).Error; err != nil {
		fmt.Println("Failed to create admin user:", err)
		return
	}

	fmt.Println("Successfully seeded admin user")
}

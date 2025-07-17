package database

import (
	"deck/enums"
	"deck/helpers"
	"deck/models"
	"fmt"
	"log"
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

func SeedProducts() {
	if DB == nil {
		fmt.Println("Error: DB is nil in SeedUser()")
		return
	}

	var products = []models.Product{
		{
			Name:        "Kopi Lampung",
			Category:    enums.Classic,
			Description: "Kopi Lampung adalah kopi robusta yang terkenal dengan cita rasa yang kuat dan aroma yang khas. Ditanam di dataran tinggi Lampung, kopi ini memiliki keasaman rendah dan body yang tebal.",
			Price:       50000,
			IsAvailable: true,
		},
		{
			Name:        "Kopi Aceh Gayo",
			Category:    enums.Classic,
			Description: "Kopi Aceh Gayo adalah kopi arabika yang ditanam di dataran tinggi Gayo, Aceh. Dikenal dengan cita rasa yang kompleks, kopi ini memiliki aroma floral dan fruity yang khas.",
			Price:       60000,
			IsAvailable: true,
		},
		{
			Name:        "Indomie Intel Telur",
			Category:    enums.Other,
			Description: "Indomie Intel Telur adalah varian mie instan yang dilengkapi dengan bumbu spesial dan telur. Cocok untuk sarapan cepat atau camilan di sore hari.",
			Price:       15000,
			IsAvailable: true,
		},
		{
			Name:        "Teh Kampleng",
			Category:    enums.Tea,
			Description: "Teh Kampleng adalah teh herbal tradisional yang terbuat dari kamplengan tangan. Dikenal dengan manfaat kesehatan yang baik, teh ini memiliki rasa yang segar dan menantang.",
			Price:       20000,
			IsAvailable: true,
		},
		{
			Name:        "Es Krim Lemon",
			Category:    enums.IceCream,
			Description: "Es Krim Lemon adalah es krim segar dengan rasa lemon yang asam manis. Cocok untuk menyegarkan hari yang panas.",
			Price:       20000,
			IsAvailable: true,
		},
	}

	for _, product := range products {
		if err := DB.Create(&product).Error; err != nil {
			log.Printf("Error creating product %s: %v", product.Name, err)
		}
	}

	log.Println("Products seeded successfully")
}

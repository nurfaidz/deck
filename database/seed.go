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

	result := DB.Where("Username = ?", "admin").FirstOrCreate(&user, models.User{
		Username: "admin",
		Email:    "admin@gmail.com",
		Password: helpers.HashPassword("password"),
	})

	if result.Error != nil {
		fmt.Println("Failed to create admin user:", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		fmt.Println("Successfully created admin user")
	} else {
		fmt.Println("Admin user already exists")
	}
}

func SeedProducts() {
	if DB == nil {
		fmt.Println("Error: DB is nil in SeedProducts()")
		return
	}

	productsToSeed := []models.Product{
		{
			Name:        "Kopi Lampung",
			Category:    enums.MainCourse,
			Description: "Kopi Lampung adalah kopi robusta yang terkenal dengan cita rasa yang kuat dan aroma yang khas. Ditanam di dataran tinggi Lampung, kopi ini memiliki keasaman rendah dan body yang tebal.",
			Price:       50000,
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
			Category:    enums.MainCourse,
			Description: "Teh Kampleng adalah teh herbal tradisional yang terbuat dari kamplengan tangan. Dikenal dengan manfaat kesehatan yang baik, teh ini memiliki rasa yang segar dan menantang.",
			Price:       20000,
			IsAvailable: true,
		},
	}

	var createdCount int
	var existingCount int

	for _, productData := range productsToSeed {
		var product models.Product

		result := DB.Where("Name = ?", productData.Name).FirstOrCreate(&product, productData)

		if result.Error != nil {
			log.Printf("Error creating/finding product '%s': %v", productData.Name, result.Error)
			continue
		}

		if result.RowsAffected > 0 {
			createdCount++
			log.Printf("Created product: %s", productData.Name)
		} else {
			existingCount++
			log.Printf("Product already exists: %s", productData.Name)
		}

		log.Printf("Product seeding completed - Created: %d, Existing: %d, Total: %d",
			createdCount, existingCount, len(productsToSeed))
	}
}

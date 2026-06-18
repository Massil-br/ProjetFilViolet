package config

import (
	"fmt"
	"log"
	"os"

	"ProjetFilViolet/backend/api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	db.AutoMigrate(
		&models.User{},
		&models.Saloon{},
		&models.Table{},
		&models.Player{},
		&models.TchatMessage{},
		&models.Friend{},
	)
	DB = db
	log.Println("✅ Connected to the database")
	seedDatabase()
}

func seedDatabase() {
	var saloonCount int64
	DB.Model(&models.Saloon{}).Count(&saloonCount)
	if saloonCount == 0 {
		saloon := models.Saloon{
			Name:        "Salon Principal",
			MinBet:      10,
			MaxBet:      200,
			PlayerCount: 7,
		}
		if err := DB.Create(&saloon).Error; err == nil {
			log.Printf("✅ Seeded default Saloon: %+v", saloon)
			
			table := models.Table{
				SaloonID:       saloon.ID,
				SlotsAvailable: 7,
			}
			if err := DB.Create(&table).Error; err == nil {
				log.Printf("✅ Seeded default Table: %+v", table)
			}
		}
	}
}

package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ .env file not found, using environment variables", err)
	}

	InitDatabase()
	fmt.Println("config file initialised")
}

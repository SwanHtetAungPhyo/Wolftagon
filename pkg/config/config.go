package config

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func LoadConfig() struct{} {
	pwd, _ := os.Getwd()
	filePath := filepath.Join(pwd, "./.env")
	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return struct{}{}
}

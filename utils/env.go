package utils

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Данные файла .env не удалось загрузить: %v", err)
	}
}

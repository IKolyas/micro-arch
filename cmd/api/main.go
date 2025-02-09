package main

import (
	"log"

	"github.com/IKolyas/otus-highload/config"
	"github.com/IKolyas/otus-highload/internal/infrastructure"
	"github.com/IKolyas/otus-highload/internal/infrastructure/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	config.Cnf = config.Config{}
	if err := config.Cnf.Load(); err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	database.DB = database.Connection{}
	if err := database.DB.Load(config.Cnf); err != nil {
		log.Fatalf("Ошибка при загрузке базы данных: %v", err)
	}

	router := infrastructure.Router()

	router.Listen(":" + config.Cnf.AppConfig.AppPort)

	log.Fatal(router)
}

package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Подключение к PostgreSQL через Docker Toolbox
	dsn := "host=192.168.99.100 user=postgres password=postgres dbname=org_structure port=5432 sslmode=disable TimeZone=UTC"

	fmt.Println("Connecting to database...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database successfully!")

	// Получаем базовую информацию о подключении
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	fmt.Println("Database ping successful!")
	fmt.Println("API server ready to start...")
}

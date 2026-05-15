package main

import (
	"fmt"
	"log"
	"net/http"

	"org-structure-api/internal/handler"
	"org-structure-api/internal/models"
	"org-structure-api/internal/repository"
	"org-structure-api/internal/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Подключение к БД
	dsn := "host=192.168.99.100 user=postgres password=postgres dbname=org_structure port=5432 sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	fmt.Println("Connected to database")

	// Миграции
	if err := db.AutoMigrate(&models.Department{}, &models.Employee{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations done")

	// Инициализация слоёв
	deptRepo := repository.NewDepartmentRepository(db)
	empRepo := repository.NewEmployeeRepository(db)

	deptService := service.NewDepartmentService(deptRepo, empRepo)
	empService := service.NewEmployeeService(empRepo, deptRepo)

	deptHandler := handler.NewDepartmentHandler(deptService)
	empHandler := handler.NewEmployeeHandler(empService)

	// Роутер
	r := mux.NewRouter()
	r.HandleFunc("/departments", deptHandler.CreateDepartment).Methods("POST")
	r.HandleFunc("/departments/{id}", deptHandler.GetDepartment).Methods("GET")
	r.HandleFunc("/departments/{id}", deptHandler.UpdateDepartment).Methods("PATCH")
	r.HandleFunc("/departments/{id}", deptHandler.DeleteDepartment).Methods("DELETE")
	r.HandleFunc("/departments/{id}/employees", empHandler.CreateEmployee).Methods("POST")

	// Запуск
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

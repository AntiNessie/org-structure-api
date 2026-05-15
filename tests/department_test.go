package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"org-structure-api/internal/handler"
	"org-structure-api/internal/models"
	"org-structure-api/internal/repository"
	"org-structure-api/internal/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Тест для создания отдела
func TestCreateDepartment(t *testing.T) {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect test DB: %v", err)
	}

	// Миграции
	db.AutoMigrate(&models.Department{}, &models.Employee{})

	// слои
	deptRepo := repository.NewDepartmentRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	deptService := service.NewDepartmentService(deptRepo, empRepo)
	deptHandler := handler.NewDepartmentHandler(deptService)

	// запрос
	body := `{"name":"Test Department"}`
	req, err := http.NewRequest("POST", "/departments", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ответ
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/departments", deptHandler.CreateDepartment).Methods("POST")
	router.ServeHTTP(rr, req)

	// Проверка статуса
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	// Проверка отдела
	var dept models.Department
	err = json.Unmarshal(rr.Body.Bytes(), &dept)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if dept.Name != "Test Department" {
		t.Errorf("Expected name 'Test Department', got '%s'", dept.Name)
	}

	if dept.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

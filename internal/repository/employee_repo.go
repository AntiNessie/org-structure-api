package repository

import (
	"org-structure-api/internal/models"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(emp *models.Employee) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(emp *models.Employee) error {
	return r.db.Create(emp).Error
}

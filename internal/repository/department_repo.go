package repository

import (
	"org-structure-api/internal/models"

	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(dept *models.Department) error
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(dept *models.Department) error {
	return r.db.Create(dept).Error
}

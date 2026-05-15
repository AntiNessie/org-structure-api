package repository

import (
	"errors"
	"org-structure-api/internal/models"

	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(dept *models.Department) error
	GetByID(id uint) (*models.Department, error)
	GetByIDWithEmployees(id uint) (*models.Department, error)
	Update(dept *models.Department) error
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

func (r *departmentRepository) GetByID(id uint) (*models.Department, error) {
	var dept models.Department
	err := r.db.First(&dept, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dept, err
}

func (r *departmentRepository) GetByIDWithEmployees(id uint) (*models.Department, error) {
	var dept models.Department
	err := r.db.Preload("Employees").First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) Update(dept *models.Department) error {
	return r.db.Save(dept).Error
}

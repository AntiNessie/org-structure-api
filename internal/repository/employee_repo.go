package repository

import (
	"org-structure-api/internal/models"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(emp *models.Employee) error
	GetByDepartmentID(deptID uint) ([]models.Employee, error) // НОВЫЙ
	DeleteByDepartmentID(deptID uint) error                   // НОВЫЙ
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

func (r *employeeRepository) GetByDepartmentID(deptID uint) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Where("department_id = ?", deptID).Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) DeleteByDepartmentID(deptID uint) error {
	return r.db.Where("department_id = ?", deptID).Delete(&models.Employee{}).Error
}

package service

import (
	"errors"
	"org-structure-api/internal/models"
	"org-structure-api/internal/repository"
	"strings"
	"time"
)

var (
	ErrEmployeeNameRequired = errors.New("full name is required and must be 1-200 characters")
	ErrPositionRequired     = errors.New("position is required and must be 1-200 characters")
)

type EmployeeService interface {
	Create(departmentID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error)
}

type employeeService struct {
	empRepo  repository.EmployeeRepository
	deptRepo repository.DepartmentRepository
}

func NewEmployeeService(empRepo repository.EmployeeRepository, deptRepo repository.DepartmentRepository) EmployeeService {
	return &employeeService{
		empRepo:  empRepo,
		deptRepo: deptRepo,
	}
}

func (s *employeeService) Create(departmentID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error) {
	// Проверяем существует ли отдел
	dept, err := s.deptRepo.GetByID(departmentID)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, ErrDepartmentNotFound
	}

	// Валидация full_name
	fullName = strings.TrimSpace(fullName)
	if len(fullName) == 0 || len(fullName) > 200 {
		return nil, ErrEmployeeNameRequired
	}

	// Валидация position
	position = strings.TrimSpace(position)
	if len(position) == 0 || len(position) > 200 {
		return nil, ErrPositionRequired
	}

	employee := &models.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
		CreatedAt:    time.Now(),
	}

	err = s.empRepo.Create(employee)
	return employee, err
}

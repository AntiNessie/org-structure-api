package service

import (
	"errors"
	"org-structure-api/internal/models"
	"org-structure-api/internal/repository"
	"strings"
)

var (
	ErrNameRequired = errors.New("name is required and must be 1-200 characters")
)

type DepartmentService interface {
	Create(name string, parentID *uint) (*models.Department, error)
}

type departmentService struct {
	deptRepo repository.DepartmentRepository
}

func NewDepartmentService(deptRepo repository.DepartmentRepository) DepartmentService {
	return &departmentService{deptRepo: deptRepo}
}

func (s *departmentService) Create(name string, parentID *uint) (*models.Department, error) {

	name = strings.TrimSpace(name)
	if len(name) == 0 || len(name) > 200 {
		return nil, ErrNameRequired
	}

	department := &models.Department{
		Name:     name,
		ParentID: parentID,
	}

	err := s.deptRepo.Create(department)
	if err != nil {
		return nil, err
	}

	return department, nil
}

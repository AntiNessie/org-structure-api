package service

import (
	"errors"
	"org-structure-api/internal/models"
	"org-structure-api/internal/repository"
	"strings"
)

var (
	ErrNameRequired       = errors.New("name is required and must be 1-200 characters")
	ErrDepartmentNotFound = errors.New("department not found")
)

type DepartmentService interface {
	Create(name string, parentID *uint) (*models.Department, error)
	GetByID(id uint, includeEmployees bool) (*models.Department, error)
	Update(id uint, name *string, parentID *uint) (*models.Department, error)
	Delete(id uint, mode string, reassignToID *uint) error
}

type departmentService struct {
	deptRepo repository.DepartmentRepository
	empRepo  repository.EmployeeRepository
}

func NewDepartmentService(deptRepo repository.DepartmentRepository, empRepo repository.EmployeeRepository) DepartmentService {
	return &departmentService{
		deptRepo: deptRepo,
		empRepo:  empRepo,
	}
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
	return department, err
}

func (s *departmentService) GetByID(id uint, includeEmployees bool) (*models.Department, error) {
	var dept *models.Department
	var err error

	if includeEmployees {
		dept, err = s.deptRepo.GetByIDWithEmployees(id)
	} else {
		dept, err = s.deptRepo.GetByID(id)
	}

	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, ErrDepartmentNotFound
	}

	return dept, nil
}

func (s *departmentService) Update(id uint, name *string, parentID *uint) (*models.Department, error) {

	dept, err := s.deptRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, ErrDepartmentNotFound
	}

	if name != nil {
		newName := strings.TrimSpace(*name)
		if len(newName) == 0 || len(newName) > 200 {
			return nil, ErrNameRequired
		}
		dept.Name = newName
	}

	if parentID != nil {
		if *parentID > 0 {
			parent, err := s.deptRepo.GetByID(*parentID)
			if err != nil {
				return nil, err
			}
			if parent == nil {
				return nil, errors.New("parent department not found")
			}
			if *parentID == id {
				return nil, errors.New("cannot set department as its own parent")
			}
		}
		dept.ParentID = parentID
	}

	err = s.deptRepo.Update(dept)
	return dept, err
}

func (s *departmentService) Delete(id uint, mode string, reassignToID *uint) error {

	dept, err := s.deptRepo.GetByID(id)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	if mode == "cascade" {
		// Удаляем всех сотрудников отдела
		if err := s.empRepo.DeleteByDepartmentID(id); err != nil {
			return err
		}
		// Удаляем отдел
		return s.deptRepo.Delete(id)
	}

	if mode == "reassign" {
		if reassignToID == nil {
			return errors.New("reassign_to_department_id is required")
		}
		targetDept, err := s.deptRepo.GetByID(*reassignToID)
		if err != nil {
			return err
		}
		if targetDept == nil {
			return errors.New("target department not found")
		}

		employees, err := s.empRepo.GetByDepartmentID(id)
		if err != nil {
			return err
		}

		// Обновляем department_id у каждого сотрудника

		for _, emp := range employees {
			emp.DepartmentID = *reassignToID
			if err := s.empRepo.Create(&emp); err != nil {
				return err
			}
		}

		// Удаляем отдел
		return s.deptRepo.Delete(id)
	}

	return errors.New("mode must be 'cascade' or 'reassign'")
}

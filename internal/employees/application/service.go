package application

import (
	"context"
	"strings"
	"time"

	"github.com/jluisv16/hcm-go/internal/employees/domain"
)

type Service struct {
	repository domain.Repository
}

type UpsertEmployeeInput struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	Department string
	Role       string
	Salary     float64
	HireDate   time.Time
}

func NewService(repository domain.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) List(ctx context.Context) ([]domain.Employee, error) {
	return s.repository.List(ctx)
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.Employee, error) {
	return s.repository.GetByID(ctx, strings.TrimSpace(id))
}

func (s *Service) Create(ctx context.Context, input UpsertEmployeeInput) (domain.Employee, error) {
	id, err := s.repository.NextID(ctx)
	if err != nil {
		return domain.Employee{}, err
	}

	employee := domain.Employee{
		ID:         id,
		FirstName:  strings.TrimSpace(input.FirstName),
		LastName:   strings.TrimSpace(input.LastName),
		Email:      strings.ToLower(strings.TrimSpace(input.Email)),
		Department: strings.TrimSpace(input.Department),
		Role:       strings.TrimSpace(input.Role),
		Salary:     input.Salary,
		HireDate:   input.HireDate,
	}

	if err := employee.Validate(); err != nil {
		return domain.Employee{}, err
	}

	inUse, err := s.repository.EmailInUse(ctx, employee.Email, "")
	if err != nil {
		return domain.Employee{}, err
	}

	if inUse {
		return domain.Employee{}, domain.ErrEmailAlreadyInUse
	}

	if err := s.repository.Create(ctx, employee); err != nil {
		return domain.Employee{}, err
	}

	return employee, nil
}

func (s *Service) Update(ctx context.Context, input UpsertEmployeeInput) (domain.Employee, error) {
	current, err := s.repository.GetByID(ctx, strings.TrimSpace(input.ID))
	if err != nil {
		return domain.Employee{}, err
	}

	current.FirstName = strings.TrimSpace(input.FirstName)
	current.LastName = strings.TrimSpace(input.LastName)
	current.Email = strings.ToLower(strings.TrimSpace(input.Email))
	current.Department = strings.TrimSpace(input.Department)
	current.Role = strings.TrimSpace(input.Role)
	current.Salary = input.Salary
	current.HireDate = input.HireDate

	if err := current.Validate(); err != nil {
		return domain.Employee{}, err
	}

	inUse, err := s.repository.EmailInUse(ctx, current.Email, current.ID)
	if err != nil {
		return domain.Employee{}, err
	}

	if inUse {
		return domain.Employee{}, domain.ErrEmailAlreadyInUse
	}

	if err := s.repository.Update(ctx, current); err != nil {
		return domain.Employee{}, err
	}

	return current, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, strings.TrimSpace(id))
}

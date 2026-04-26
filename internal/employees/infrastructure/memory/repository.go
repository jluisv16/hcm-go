package memory

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/jluisv16/hcm-go/internal/employees/domain"
)

type Repository struct {
	mu           sync.RWMutex
	employees    map[string]domain.Employee
	nextSequence int
}

func NewRepository(seed []domain.Employee) *Repository {
	repository := &Repository{
		employees: make(map[string]domain.Employee, len(seed)),
	}

	maxSequence := 0
	for _, employee := range seed {
		repository.employees[employee.ID] = employee
		sequence := parseEmployeeSequence(employee.ID)
		if sequence > maxSequence {
			maxSequence = sequence
		}
	}

	repository.nextSequence = maxSequence

	return repository
}

func (r *Repository) List(_ context.Context) ([]domain.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.Employee, 0, len(r.employees))
	for _, employee := range r.employees {
		list = append(list, employee)
	}

	sort.Slice(list, func(i, j int) bool {
		return parseEmployeeSequence(list[i].ID) < parseEmployeeSequence(list[j].ID)
	})

	return list, nil
}

func (r *Repository) GetByID(_ context.Context, id string) (domain.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	employee, exists := r.employees[id]
	if !exists {
		return domain.Employee{}, domain.ErrEmployeeNotFound
	}

	return employee, nil
}

func (r *Repository) Create(_ context.Context, employee domain.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.employees[employee.ID]; exists {
		return fmt.Errorf("employee %s already exists", employee.ID)
	}

	r.employees[employee.ID] = employee

	return nil
}

func (r *Repository) Update(_ context.Context, employee domain.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.employees[employee.ID]; !exists {
		return domain.ErrEmployeeNotFound
	}

	r.employees[employee.ID] = employee

	return nil
}

func (r *Repository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.employees[id]; !exists {
		return domain.ErrEmployeeNotFound
	}

	delete(r.employees, id)
	return nil
}

func (r *Repository) NextID(_ context.Context) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextSequence++
	return fmt.Sprintf("EMP-%03d", r.nextSequence), nil
}

func (r *Repository) EmailInUse(_ context.Context, email, excludingID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	for _, employee := range r.employees {
		if strings.EqualFold(employee.ID, excludingID) {
			continue
		}

		if strings.EqualFold(employee.Email, normalizedEmail) {
			return true, nil
		}
	}

	return false, nil
}

func parseEmployeeSequence(id string) int {
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return 0
	}

	sequence, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}

	return sequence
}

package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrEmployeeNotFound  = errors.New("employee not found")
	ErrEmailAlreadyInUse = errors.New("employee email already in use")
)

type Employee struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	Department string
	Role       string
	Salary     float64
	HireDate   time.Time
}

func (e Employee) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("id is required")
	}

	if strings.TrimSpace(e.FirstName) == "" {
		return errors.New("first name is required")
	}

	if strings.TrimSpace(e.LastName) == "" {
		return errors.New("last name is required")
	}

	if strings.TrimSpace(e.Email) == "" {
		return errors.New("email is required")
	}

	if !strings.Contains(e.Email, "@") {
		return fmt.Errorf("email %q is invalid", e.Email)
	}

	if strings.TrimSpace(e.Department) == "" {
		return errors.New("department is required")
	}

	if strings.TrimSpace(e.Role) == "" {
		return errors.New("role is required")
	}

	if e.Salary <= 0 {
		return errors.New("salary must be greater than zero")
	}

	if e.HireDate.IsZero() {
		return errors.New("hire date is required")
	}

	return nil
}

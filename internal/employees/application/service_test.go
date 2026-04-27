package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jluisv16/hcm-go/internal/employees/domain"
	"github.com/jluisv16/hcm-go/internal/employees/infrastructure/memory"
)

func newTestService() *Service {
	return NewService(memory.NewRepository(memory.SeedEmployees()))
}

func TestCreateNormalizesInput(t *testing.T) {
	t.Parallel()

	service := newTestService()
	ctx := context.Background()

	input := UpsertEmployeeInput{
		FirstName:  "  John  ",
		LastName:   "  Doe ",
		Email:      "  JOHN.DOE@HCM.LOCAL ",
		Department: " IT ",
		Role:       " Platform Engineer ",
		Salary:     3400,
		HireDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	created, err := service.Create(ctx, input)
	if err != nil {
		t.Fatalf("expected create to succeed, got error: %v", err)
	}

	if created.ID != "EMP-011" {
		t.Fatalf("expected generated id EMP-011, got %s", created.ID)
	}

	if created.FirstName != "John" || created.LastName != "Doe" {
		t.Fatalf("expected names to be trimmed, got %q %q", created.FirstName, created.LastName)
	}

	if created.Email != "john.doe@hcm.local" {
		t.Fatalf("expected email normalized to lowercase, got %s", created.Email)
	}
}

func TestCreateDuplicateEmailReturnsDomainError(t *testing.T) {
	t.Parallel()

	service := newTestService()
	ctx := context.Background()

	_, err := service.Create(ctx, UpsertEmployeeInput{
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "ana.lopez@hcm.local",
		Department: "IT",
		Role:       "SRE",
		Salary:     3000,
		HireDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrEmailAlreadyInUse) {
		t.Fatalf("expected ErrEmailAlreadyInUse, got %v", err)
	}
}

func TestUpdateNotFoundReturnsDomainError(t *testing.T) {
	t.Parallel()

	service := newTestService()
	ctx := context.Background()

	_, err := service.Update(ctx, UpsertEmployeeInput{
		ID:         "EMP-999",
		FirstName:  "Missing",
		LastName:   "User",
		Email:      "missing@hcm.local",
		Department: "IT",
		Role:       "QA",
		Salary:     2000,
		HireDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrEmployeeNotFound) {
		t.Fatalf("expected ErrEmployeeNotFound, got %v", err)
	}
}

func TestUpdateDuplicateEmailReturnsDomainError(t *testing.T) {
	t.Parallel()

	service := newTestService()
	ctx := context.Background()

	_, err := service.Update(ctx, UpsertEmployeeInput{
		ID:         "EMP-001",
		FirstName:  "Luis",
		LastName:   "Vargas",
		Email:      "ana.lopez@hcm.local",
		Department: "HR",
		Role:       "HR Analyst",
		Salary:     1900,
		HireDate:   time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrEmailAlreadyInUse) {
		t.Fatalf("expected ErrEmailAlreadyInUse, got %v", err)
	}
}

func TestDeleteRemovesEmployee(t *testing.T) {
	t.Parallel()

	service := newTestService()
	ctx := context.Background()

	if err := service.Delete(ctx, " EMP-001 "); err != nil {
		t.Fatalf("expected delete to succeed, got error: %v", err)
	}

	_, err := service.GetByID(ctx, "EMP-001")
	if !errors.Is(err, domain.ErrEmployeeNotFound) {
		t.Fatalf("expected ErrEmployeeNotFound after delete, got %v", err)
	}
}


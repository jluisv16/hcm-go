package memory

import (
	"time"

	"github.com/jluisv16/hcm-go/internal/employees/domain"
)

func SeedEmployees() []domain.Employee {
	return []domain.Employee{
		{
			ID:         "EMP-001",
			FirstName:  "Luis",
			LastName:   "Vargas",
			Email:      "luis.vargas@hcm.local",
			Department: "HR",
			Role:       "HR Analyst",
			Salary:     1800.00,
			HireDate:   time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-002",
			FirstName:  "Ana",
			LastName:   "Lopez",
			Email:      "ana.lopez@hcm.local",
			Department: "IT",
			Role:       "Backend Engineer",
			Salary:     3200.00,
			HireDate:   time.Date(2020, 6, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-003",
			FirstName:  "Carlos",
			LastName:   "Mendez",
			Email:      "carlos.mendez@hcm.local",
			Department: "Finance",
			Role:       "Accountant",
			Salary:     2200.00,
			HireDate:   time.Date(2019, 8, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-004",
			FirstName:  "Mariana",
			LastName:   "Rojas",
			Email:      "mariana.rojas@hcm.local",
			Department: "Operations",
			Role:       "Operations Lead",
			Salary:     2900.00,
			HireDate:   time.Date(2022, 2, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-005",
			FirstName:  "Pedro",
			LastName:   "Silva",
			Email:      "pedro.silva@hcm.local",
			Department: "IT",
			Role:       "DevOps Engineer",
			Salary:     3500.00,
			HireDate:   time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-006",
			FirstName:  "Gabriela",
			LastName:   "Torres",
			Email:      "gabriela.torres@hcm.local",
			Department: "Marketing",
			Role:       "Marketing Specialist",
			Salary:     2100.00,
			HireDate:   time.Date(2021, 9, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-007",
			FirstName:  "Javier",
			LastName:   "Castro",
			Email:      "javier.castro@hcm.local",
			Department: "HR",
			Role:       "Recruiter",
			Salary:     1700.00,
			HireDate:   time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-008",
			FirstName:  "Sofia",
			LastName:   "Perez",
			Email:      "sofia.perez@hcm.local",
			Department: "Finance",
			Role:       "Financial Analyst",
			Salary:     2400.00,
			HireDate:   time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-009",
			FirstName:  "Diego",
			LastName:   "Ramirez",
			Email:      "diego.ramirez@hcm.local",
			Department: "Operations",
			Role:       "Process Analyst",
			Salary:     2600.00,
			HireDate:   time.Date(2019, 4, 11, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         "EMP-010",
			FirstName:  "Valeria",
			LastName:   "Ortega",
			Email:      "valeria.ortega@hcm.local",
			Department: "IT",
			Role:       "QA Engineer",
			Salary:     2300.00,
			HireDate:   time.Date(2022, 7, 21, 0, 0, 0, 0, time.UTC),
		},
	}
}

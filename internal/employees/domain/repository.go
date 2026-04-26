package domain

import "context"

type Repository interface {
	List(ctx context.Context) ([]Employee, error)
	GetByID(ctx context.Context, id string) (Employee, error)
	Create(ctx context.Context, employee Employee) error
	Update(ctx context.Context, employee Employee) error
	Delete(ctx context.Context, id string) error
	NextID(ctx context.Context) (string, error)
	EmailInUse(ctx context.Context, email, excludingID string) (bool, error)
}

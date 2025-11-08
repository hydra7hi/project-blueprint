package database

import (
	"context"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// SQLClientInterface
type SQLClientInterface interface {
	CreateUser(ctx context.Context, name, email string, age int32) (*UserRow, error)
	GetUser(ctx context.Context, id int) (*UserRow, error)
	UpdateUser(ctx context.Context, id int, name, email string, age int32) (*UserRow, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, limit, offset int) ([]*UserRow, error)
	CountUsers(ctx context.Context) (int, error)
}

package database

import (
	"fmt"
	"time"

	pb "grpc-services/user/proto"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// UserRow
// Represent a single DB row.
type UserRow struct {
	ID        int
	Name      string
	Email     string
	Age       int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToProto
func (u *UserRow) ToProto() *pb.User {
	return &pb.User{
		Id:        fmt.Sprintf("%d", u.ID),
		Name:      u.Name,
		Email:     u.Email,
		Age:       u.Age,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

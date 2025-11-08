package server

import (
	"context"
	"strings"

	pb "grpc-services/user/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handlers
// Ensure the request is valid.
// Before preforming the logic.
//
// Returns:
//   - Proper response.
//
// Errors:
//   - InvalidArgument: With message indicating the validation issue.
//   - Other: the returned logic error.

// CreateUser handler
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	// Validate Request
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	}
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email cannot be empty")
	}
	if req.GetAge() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "age must be positive")
	}
	if !strings.Contains(req.GetEmail(), "@") {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	// Execute Logic
	return s.createUser(ctx, req)
}

// GetUser handler
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	// Validate Request
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID cannot be empty")
	}

	// Execute Logic
	return s.getUser(ctx, req)
}

// UpdateUser handler
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	// Validate Request
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID cannot be empty")
	}
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	}
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email cannot be empty")
	}
	if req.GetAge() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "age must be positive")
	}
	if !strings.Contains(req.GetEmail(), "@") {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	// Execute Logic
	return s.updateUser(ctx, req)
}

// DeleteUser handler
func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Validate Request
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID cannot be empty")
	}

	// Execute Logic
	return s.deleteUser(ctx, req)
}

// ListUsers handler
func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Validate Request
	if req.GetPage() < 0 {
		return nil, status.Error(codes.InvalidArgument, "page cannot be negative")
	}
	if req.GetLimit() < 0 {
		return nil, status.Error(codes.InvalidArgument, "limit cannot be negative")
	}

	// Execute Logic
	return s.listUsers(ctx, req)
}

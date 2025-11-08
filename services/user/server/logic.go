package server

import (
	"context"
	"fmt"
	"strconv"

	pb "grpc-services/user/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Logic
// Preforms the logic behind and rpc.
// Assumes all given values are validated before calling the function.
//
// Returns:
//   - Proper response.
//
// For simplicity, the code is done in a way that
//
// (Note: A real implemention would split further split the errors into:
// internal, Unavailable, NotFound based on db error type)

// createUser
//
// Errors:
//   - Internal: When failing to create user in DB.
//
// (Note: A real implemention would split the errors into: internal, Unavailable, NotFound based on error type)
func (s *Server) createUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := s.DB.CreateUser(ctx, req.GetName(), req.GetEmail(), req.GetAge())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create user: %v", err))
	}
	return &pb.UserResponse{User: user.ToProto()}, nil
}

// getUser
//
// Errors:
//   - NotFound: When failing to find user in DB.
//
// (Note: A real implemention would split the errors into: internal, Unavailable, NotFound based on error type)
func (s *Server) getUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, err
	}

	user, err := s.DB.GetUser(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &pb.UserResponse{User: user.ToProto()}, nil
}

// updateUser
//
// Errors:
//   - NotFound: When failing to find user in DB.
//
// (Note: A real implemention would split the errors into: internal, Unavailable, NotFound based on error type)
func (s *Server) updateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, err
	}

	user, err := s.DB.UpdateUser(ctx, id, req.GetName(), req.GetEmail(), req.GetAge())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &pb.UserResponse{User: user.ToProto()}, nil
}

// deleteUser
//
// Errors:
//   - NotFound: When failing to find user in DB.
//
// (Note: A real implemention would split the errors into: internal, Unavailable, NotFound based on error type)
func (s *Server) deleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, err
	}

	err = s.DB.DeleteUser(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

// listUsers
// Return a list of users with paging
// Limits page size to 100
// defaults (page, limit) to: (1, 10) if missing.
//
// Errors:
//   - Internal: When failing to find user in DB.
//
// (Note: A real implemention would split the errors into: internal, Unavailable, NotFound based on error type)
func (s *Server) listUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	page := int(req.GetPage())
	if page < 1 {
		page = 1
	}

	limit := int(req.GetLimit())
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, err := s.DB.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list users: %v", err))
	}

	total, err := s.DB.CountUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to count users: %v", err))
	}

	var protoUsers []*pb.User
	for _, user := range users {
		protoUsers = append(protoUsers, user.ToProto())
	}

	return &pb.ListUsersResponse{
		Users: protoUsers,
		Total: int32(total),
		Page:  int32(page),
		Limit: int32(limit),
	}, nil
}

// parseID
// parse a given idString
// Returns:
//   - id value as an int
//
// Errors:
//   - InvalidArgument: When failing to find user in DB.
func parseID(idString string) (int, error) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, status.Error(codes.InvalidArgument, "invalid user ID")
	}
	return id, nil
}

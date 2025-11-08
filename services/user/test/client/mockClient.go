package client

import (
	"context"

	pb "grpc-services/user/proto"

	"google.golang.org/grpc"
)

type GRPCClientInterface interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error)
	GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error)
	UpdateUser(ctx context.Context, in *pb.UpdateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error)
	DeleteUser(ctx context.Context, in *pb.DeleteUserRequest, opts ...grpc.CallOption) (*pb.DeleteUserResponse, error)
	ListUsers(ctx context.Context, in *pb.ListUsersRequest, opts ...grpc.CallOption) (*pb.ListUsersResponse, error)
}

type MockGRPCClient struct {
	// Responses
	CreateUserResponse *pb.UserResponse
	GetUserResponse    *pb.UserResponse
	UpdateUserResponse *pb.UserResponse
	DeleteUserResponse *pb.DeleteUserResponse
	ListUsersResponse  *pb.ListUsersResponse

	// Errors
	CreateUserError error
	GetUserError    error
	UpdateUserError error
	DeleteUserError error
	ListUsersError  error

	// Call counts
	CreateUserCount int
	GetUserCount    int
	UpdateUserCount int
	DeleteUserCount int
	ListUsersCount  int

	// Last requests
	LastCreateUserRequest *pb.CreateUserRequest
	LastGetUserRequest    *pb.GetUserRequest
	LastUpdateUserRequest *pb.UpdateUserRequest
	LastDeleteUserRequest *pb.DeleteUserRequest
	LastListUsersRequest  *pb.ListUsersRequest
}

func (c *MockGRPCClient) CreateUser(ctx context.Context, in *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	c.CreateUserCount++
	c.LastCreateUserRequest = in
	if c.CreateUserError != nil {
		return nil, c.CreateUserError
	}
	return c.CreateUserResponse, nil
}

func (c *MockGRPCClient) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	c.GetUserCount++
	c.LastGetUserRequest = in
	if c.GetUserError != nil {
		return nil, c.GetUserError
	}
	return c.GetUserResponse, nil
}

func (c *MockGRPCClient) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	c.UpdateUserCount++
	c.LastUpdateUserRequest = in
	if c.UpdateUserError != nil {
		return nil, c.UpdateUserError
	}
	return c.UpdateUserResponse, nil
}

func (c *MockGRPCClient) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest, opts ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
	c.DeleteUserCount++
	c.LastDeleteUserRequest = in
	if c.DeleteUserError != nil {
		return nil, c.DeleteUserError
	}
	return c.DeleteUserResponse, nil
}

func (c *MockGRPCClient) ListUsers(ctx context.Context, in *pb.ListUsersRequest, opts ...grpc.CallOption) (*pb.ListUsersResponse, error) {
	c.ListUsersCount++
	c.LastListUsersRequest = in
	if c.ListUsersError != nil {
		return nil, c.ListUsersError
	}
	return c.ListUsersResponse, nil
}

// Helper methods for test setup
func (c *MockGRPCClient) Reset() {
	*c = MockGRPCClient{}
}

func (c *MockGRPCClient) SetCreateUserResponse(user *pb.User) {
	c.CreateUserResponse = &pb.UserResponse{User: user}
}

func (c *MockGRPCClient) SetGetUserResponse(user *pb.User) {
	c.GetUserResponse = &pb.UserResponse{User: user}
}

func (c *MockGRPCClient) SetUpdateUserResponse(user *pb.User) {
	c.UpdateUserResponse = &pb.UserResponse{User: user}
}

func (c *MockGRPCClient) SetDeleteUserResponse(success bool) {
	c.DeleteUserResponse = &pb.DeleteUserResponse{Success: success}
}

func (c *MockGRPCClient) SetListUsersResponse(users []*pb.User, total int32) {
	c.ListUsersResponse = &pb.ListUsersResponse{
		Users: users,
		Total: total,
		Page:  1,
		Limit: int32(len(users)),
	}
}

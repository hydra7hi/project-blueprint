package client

import (
	"context"

	"google.golang.org/grpc"

	pb "grpc-services/user/proto"
)

// GRPCClientInterface
// To be included and used in other servies.
// Makes mocking other services easy for unit tests.
type GRPCClientInterface interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error)
	GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error)
	UpdateUser(ctx context.Context, in *pb.UpdateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error)
	DeleteUser(ctx context.Context, in *pb.DeleteUserRequest, opts ...grpc.CallOption) (*pb.DeleteUserResponse, error)
	ListUsers(ctx context.Context, in *pb.ListUsersRequest, opts ...grpc.CallOption) (*pb.ListUsersResponse, error)
}

package client

import (
	"context"

	"google.golang.org/grpc"

	pb "grpc-services/user/proto"
)

// GRPCClient implements GRPCClientInterface
// Can be used in other services.
type GRPCClient struct {
	Client pb.UserServiceClient
	conn   *grpc.ClientConn
	Cancel context.CancelFunc
}

// NewGRPCClient
// Creates new GRPCClient.
//
// Returns:
//   - GRPCClient
//
// Errors:
//   - If it fails to Dial the grpc connection.
func NewGRPCClient(ctx context.Context, serverAddress string) (*GRPCClient, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	_, cancel := context.WithCancel(ctx)

	client := pb.NewUserServiceClient(conn)
	return &GRPCClient{
		Client: client,
		conn:   conn,
		Cancel: cancel,
	}, nil
}

func (c *GRPCClient) Close() {
	if c.Cancel != nil {
		c.Cancel()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *GRPCClient) CreateUser(ctx context.Context, in *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	return c.Client.CreateUser(ctx, in, opts...)
}

func (c *GRPCClient) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	return c.Client.GetUser(ctx, in, opts...)
}

func (c *GRPCClient) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	return c.Client.UpdateUser(ctx, in, opts...)
}

func (c *GRPCClient) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest, opts ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
	return c.Client.DeleteUser(ctx, in, opts...)
}

func (c *GRPCClient) ListUsers(ctx context.Context, in *pb.ListUsersRequest, opts ...grpc.CallOption) (*pb.ListUsersResponse, error) {
	return c.Client.ListUsers(ctx, in, opts...)
}

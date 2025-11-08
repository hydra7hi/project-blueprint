package client

import (
	"context"
	"fmt"

	pb "grpc-services/operation/proto"

	"google.golang.org/grpc"
)

// GRPCClient
// Wraps the gRPC client connection and service client.
type GRPCClient struct {
	Conn   *grpc.ClientConn
	Client pb.OperationServiceClient
}

// NewGRPCClient
// Creates a new gRPC client connection to the operation service.
//
// Returns:
//   - *GRPCClient: Client instance with connection and service client
//
// Errors:
//   - Failed to establish gRPC connection
func NewGRPCClient(ctx context.Context, address string) (*GRPCClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to operation service: %v", err)
	}

	client := pb.NewOperationServiceClient(conn)

	return &GRPCClient{
		Conn:   conn,
		Client: client,
	}, nil
}

// Close
// Closes the gRPC client connection.
// Should be called when client is no longer needed.
func (c *GRPCClient) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

func (c *GRPCClient) StartOperation(ctx context.Context, in *pb.StartOperationRequest, opts ...grpc.CallOption) (*pb.StartOperationResponse, error) {
	return c.Client.StartOperation(ctx, in, opts...)
}

func (c *GRPCClient) CheckProcess(ctx context.Context, in *pb.CheckProcessRequest, opts ...grpc.CallOption) (*pb.CheckProcessResponse, error) {
	return c.Client.CheckProcess(ctx, in, opts...)
}

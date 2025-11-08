package client

import (
	"context"

	"google.golang.org/grpc"

	pb "grpc-services/operation/proto"
)

// GRPCClientInterface
// To be included and used in other services.
// Makes mocking other services easy for unit tests.
type GRPCClientInterface interface {
	StartOperation(ctx context.Context, in *pb.StartOperationRequest, opts ...grpc.CallOption) (*pb.StartOperationResponse, error)
	CheckProcess(ctx context.Context, in *pb.CheckProcessRequest, opts ...grpc.CallOption) (*pb.CheckProcessResponse, error)
}

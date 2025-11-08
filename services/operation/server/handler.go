package server

import (
	"context"

	pb "grpc-services/operation/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StartOperation
// Queues a new operation for background processing.
// Returns operation ID for status tracking.
//
// Returns:
//   - StartOperationResponse with operation ID
//
// Errors:
//   - InvalidArgument: Operation data is invalid
//   - Internal: Failed to queue operation
func (s *Server) StartOperation(ctx context.Context, req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	// Validate Request
	if req.GetOperationData() == nil {
		return nil, status.Error(codes.InvalidArgument, "operation data cannot be empty")
	}

	// Execute Logic
	return s.startOperation(ctx, req)
}

// CheckProcess
// Retrieves current status of a long-running operation.
// Used to track progress and completion.
//
// Returns:
//   - CheckProcessResponse with current operation state
//
// Errors:
//   - InvalidArgument: Operation ID is empty
//   - NotFound: Operation ID does not exist
//   - Internal: Failed to retrieve operation status
func (s *Server) CheckProcess(ctx context.Context, req *pb.CheckProcessRequest) (*pb.CheckProcessResponse, error) {
	if req.GetOperationId() == "" {
		return s.getLatestOperation(ctx, req)
	}

	// Execute Logic
	return s.checkProcess(ctx, req)
}

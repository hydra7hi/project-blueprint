package server

import (
	"context"
	"fmt"
	"time"

	"grpc-services/operation/database"
	pb "grpc-services/operation/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Logic
// Preforms the logic behind an rpc.
// Assumes all given values are validated before calling the function.
//
// Returns:
//   - Proper response.
//
// For simplicity, the code is done in a way that
//
// (Note: A real implementation would split further split the errors into:
// internal, Unavailable, NotFound based on db error type)

// startOperation
// Creates a new operation and queues it for background processing.
//
// Errors:
//   - Internal: When failing to create operation in DB.
func (s *Server) startOperation(ctx context.Context, req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	// Generate operation ID
	opID := generateOperationID()

	// Marshal request data for storage
	marshalledReq, err := marshalOperationData(req.OperationData)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to marshal operation data: %v", err))
	}

	// Create operation in database
	operation := &database.Operation{
		ID:                opID,
		MarshalledRequest: marshalledReq,
		StepID:            int(StepInitial),
		State:             database.StatePending,
	}

	if err := s.DB.CreateOperation(ctx, operation); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create operation: %v", err))
	}

	// Start background processing
	go func() {
		// Use background context for the operation processing
		bgCtx := context.Background()
		if err := s.Processor.ProcessOperation(bgCtx, opID); err != nil {
			fmt.Printf("Operation %s failed: %v\n", opID, err)

			// Set to failed in the DB.
			if err := s.DB.UpdateOperationState(bgCtx, opID, database.StateFailed); err != nil {
				fmt.Printf("Operation %s failed: %v\n", opID, err)
			}
		} else {
			// Set to completed in the DB.
			if err := s.DB.UpdateOperationState(bgCtx, opID, database.StateCompleted); err != nil {
				fmt.Printf("Operation %s failed: %v\n", opID, err)
			}
		}

	}()

	return &pb.StartOperationResponse{
		OperationId: opID,
	}, nil
}

// checkProcess
// Retrieves the current status of an operation.
//
// Errors:
//   - NotFound: When failing to find operation in DB.
func (s *Server) checkProcess(ctx context.Context, req *pb.CheckProcessRequest) (*pb.CheckProcessResponse, error) {
	operation, err := s.DB.GetOperation(ctx, req.OperationId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "operation not found")
	}

	return &pb.CheckProcessResponse{
		OperationId: operation.ID,
		CurrentStep: int32(operation.StepID),
		TotalSteps:  int32(TotalSteps),
		State:       operation.State.String(),
		Completed:   isOperationCompleted(operation.State),
	}, nil
}

// getLatestOperation
// Retrieves the most recently created operation.
//
// Errors:
//   - NotFound: When no operations exist in DB.
func (s *Server) getLatestOperation(ctx context.Context, _ *pb.CheckProcessRequest) (*pb.CheckProcessResponse, error) {
	operation, err := s.DB.GetLatestOperation(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no operations found")
	}

	return &pb.CheckProcessResponse{
		OperationId: operation.ID,
		CurrentStep: int32(operation.StepID),
		TotalSteps:  int32(TotalSteps),
		State:       operation.State.String(),
		Completed:   isOperationCompleted(operation.State),
	}, nil
}

// generateOperationID
// Generates a unique operation ID.
// In production, use UUID or other unique identifier.
func generateOperationID() string {
	return fmt.Sprintf("op-%d", time.Now().UnixNano())
}

// marshalOperationData
// Marshals operation data for storage.
func marshalOperationData(data *pb.OperationData) ([]byte, error) {
	// Since OperationData is empty for now, return empty JSON
	// In future, this will marshal the actual operation data
	return []byte("{}"), nil
}

// isOperationCompleted
// Checks if operation has reached terminal state.
func isOperationCompleted(state database.OperationState) bool {
	return state == database.StateCompleted ||
		state == database.StateFailed ||
		state == database.StateCancelled
}

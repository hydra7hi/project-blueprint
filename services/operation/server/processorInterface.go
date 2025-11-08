package server

import (
	"context"
	"grpc-services/operation/database"
)

// OperationProcessorInterface defines the contract for operation processing
type OperationProcessorInterface interface {
	// ProcessOperation executes the operation steps in the background
	ProcessOperation(ctx context.Context, operationID string) error

	// processStepInitial starts the operation processing
	processStepInitial(ctx context.Context, operation *database.Operation) error

	// processStepListUsers lists all existing users
	processStepListUsers(ctx context.Context, operation *database.Operation) error

	// processStepDeleteUsers deletes all existing users
	processStepDeleteUsers(ctx context.Context, operation *database.Operation) error

	// processStepCreateUsers creates 5 new users
	processStepCreateUsers(ctx context.Context, operation *database.Operation) error

	// // StartBackgroundProcessor starts the background worker that processes operations
	// StartBackgroundProcessor(ctx context.Context)

	// // processPendingOperations finds and processes pending operations
	// processPendingOperations(ctx context.Context)
}

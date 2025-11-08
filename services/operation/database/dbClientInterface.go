package database

import "context"

// DBClientInterface defines the contract for database operations
type DBClientInterface interface {
	// CreateOperation creates a new operation in the database
	CreateOperation(ctx context.Context, op *Operation) error

	// GetOperation retrieves an operation by ID
	GetOperation(ctx context.Context, id string) (*Operation, error)

	// UpdateOperation updates an existing operation
	UpdateOperation(ctx context.Context, op *Operation) error

	// UpdateOperationState updates only the state of an operation
	UpdateOperationState(ctx context.Context, id string, state OperationState) error

	// UpdateOperationStep updates the step and state of an operation
	UpdateOperationStep(ctx context.Context, id string, stepID int, state OperationState) error

	// GetLatestOperation retrieves the most recently created operation
	GetLatestOperation(ctx context.Context) (*Operation, error)
}

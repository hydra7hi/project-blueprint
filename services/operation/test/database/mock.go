package database

import (
	"context"
	"fmt"
	"grpc-services/operation/database"
	"time"
)

// Mock Client
type MockClient struct {
	givenCreateError error
	givenGetError    error
	givenUpdateError error

	Operations map[string]*database.Operation
}

func NewMockClient(givenCreateError error, givenGetError error, givenUpdateError error) *MockClient {
	return &MockClient{
		givenCreateError: givenCreateError,
		givenGetError:    givenGetError,
		givenUpdateError: givenUpdateError,
		Operations:       make(map[string]*database.Operation),
	}
}

func (m *MockClient) CreateOperation(ctx context.Context, op *database.Operation) error {
	if m.givenCreateError != nil {
		return m.givenCreateError
	}

	m.Operations[op.ID] = op
	return nil
}

func (m *MockClient) GetOperation(ctx context.Context, id string) (*database.Operation, error) {
	if m.givenGetError != nil {
		return nil, m.givenGetError
	}

	op, exists := m.Operations[id]
	if !exists {
		return nil, fmt.Errorf("operation not found")
	}
	return op, nil
}

func (m *MockClient) UpdateOperation(ctx context.Context, op *database.Operation) error {
	if m.givenUpdateError != nil {
		return m.givenUpdateError
	}

	if _, exists := m.Operations[op.ID]; !exists {
		return fmt.Errorf("operation not found")
	}

	m.Operations[op.ID] = op
	return nil
}

func (m *MockClient) UpdateOperationState(ctx context.Context, id string, state database.OperationState) error {
	if m.givenUpdateError != nil {
		return m.givenUpdateError
	}

	op, exists := m.Operations[id]
	if !exists {
		return fmt.Errorf("operation not found")
	}

	op.State = state
	op.UpdatedAt = time.Now()
	return nil
}

func (m *MockClient) UpdateOperationStep(ctx context.Context, id string, stepID int, state database.OperationState) error {
	if m.givenUpdateError != nil {
		return m.givenUpdateError
	}

	op, exists := m.Operations[id]
	if !exists {
		return fmt.Errorf("operation not found")
	}

	op.StepID = stepID
	op.State = state
	op.UpdatedAt = time.Now()
	return nil
}

func (m *MockClient) GetLatestOperation(ctx context.Context) (*database.Operation, error) {
	if m.givenGetError != nil {
		return nil, m.givenGetError
	}

	var latest *database.Operation
	for _, op := range m.Operations {
		if latest == nil || op.CreatedAt.After(latest.CreatedAt) {
			latest = op
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("no operations found")
	}
	return latest, nil
}

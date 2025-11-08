package server

import (
	"context"
	"grpc-services/operation/database"
)

// Mock Processor
type MockProcessor struct {
	givenProcessError error

	ProcessedOperations []string
}

func NewMockProcessor(givenProcessError error) *MockProcessor {
	return &MockProcessor{
		givenProcessError: givenProcessError,
	}
}

func (m *MockProcessor) ProcessOperation(ctx context.Context, operationID string) error {
	if m.givenProcessError != nil {
		return m.givenProcessError
	}

	m.ProcessedOperations = append(m.ProcessedOperations, operationID)
	return nil
}

func (m *MockProcessor) processStepInitial(ctx context.Context, operation *database.Operation) error {
	return m.ProcessOperation(ctx, operation.ID)
}

func (m *MockProcessor) processStepListUsers(ctx context.Context, operation *database.Operation) error {
	return m.ProcessOperation(ctx, operation.ID)
}

func (m *MockProcessor) processStepDeleteUsers(ctx context.Context, operation *database.Operation) error {
	return m.ProcessOperation(ctx, operation.ID)
}

func (m *MockProcessor) processStepCreateUsers(ctx context.Context, operation *database.Operation) error {
	return m.ProcessOperation(ctx, operation.ID)
}

func (m *MockProcessor) StartBackgroundProcessor(ctx context.Context) {
	// Mock implementation - no background processing
}

func (m *MockProcessor) processPendingOperations(ctx context.Context) {
	// Mock implementation - no pending operations processing
}

package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc-services/operation/database"
	userCl "grpc-services/user/client"
	userpb "grpc-services/user/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OperationStep represents the steps in the operation process
type OperationStep int

const (
	// StepInitial represents the initial step before processing begins
	StepInitial OperationStep = iota
	// StepListUsers represents listing existing users
	StepListUsers
	// StepDeleteUsers represents deleting all existing users
	StepDeleteUsers
	// StepCreateUsers represents creating new users
	StepCreateUsers
	// StepCompleted represents the final completed step
	StepCompleted
)

// String returns the string representation of the operation step
func (s OperationStep) String() string {
	return [...]string{"INITIAL", "LIST_USERS", "DELETE_USERS", "CREATE_USERS", "COMPLETED"}[s]
}

// TotalSteps returns the total number of processing steps (excluding initial and completed)
const TotalSteps = StepCompleted

// OperationProcessor handles the background processing of operations
type OperationProcessor struct {
	delayExecutuion bool
	dbClient        database.DBClientInterface
	userClient      userCl.GRPCClientInterface
}

// NewOperationProcessor creates a new operation processor
func NewOperationProcessor(dbClient database.DBClientInterface, userClient userpb.UserServiceClient) *OperationProcessor {
	return &OperationProcessor{
		delayExecutuion: true, // Set to true for demo purpose when running the service.
		dbClient:        dbClient,
		userClient:      userClient,
	}
}

// NewTestOperationProcessor creates a new test operation processor
func NewTestOperationProcessor(dbClient database.DBClientInterface, userClient userpb.UserServiceClient) *OperationProcessor {
	return &OperationProcessor{
		delayExecutuion: false, // Set to false for faster execution of tests
		dbClient:        dbClient,
		userClient:      userClient,
	}
}

// ProcessOperation executes the operation steps in the background
// Handles resumption from any step in case of service restart
func (p *OperationProcessor) ProcessOperation(ctx context.Context, operationID string) error {
	// Get the current operation state from database
	operation, err := p.dbClient.GetOperation(ctx, operationID)
	if err != nil {
		return fmt.Errorf("failed to get operation: %v", err)
	}

	if p.delayExecutuion {
		// Small sleep for demo.
		// To allow for monitoring the step progress in logs.
		time.Sleep(3 * time.Second)
	}
	log.Printf("Operation [%s] :Starting Step: %s", operationID, OperationStep(operation.StepID))

	// Process based on current step
	switch OperationStep(operation.StepID) {
	case StepInitial:
		return p.processStepInitial(ctx, operation)
	case StepListUsers:
		return p.processStepListUsers(ctx, operation)
	case StepDeleteUsers:
		return p.processStepDeleteUsers(ctx, operation)
	case StepCreateUsers:
		return p.processStepCreateUsers(ctx, operation)
	default:
		return fmt.Errorf("unknown step ID: %d", operation.StepID)
	}
}

// processStepInitial starts the operation processing
func (p *OperationProcessor) processStepInitial(ctx context.Context, operation *database.Operation) error {
	log.Printf("Starting operation [%s]", operation.ID)

	// Update to first processing step
	if err := p.dbClient.UpdateOperationStep(ctx, operation.ID, int(StepListUsers), database.StateRunning); err != nil {
		return fmt.Errorf("failed to update operation step: %v", err)
	}

	// Continue processing
	return p.ProcessOperation(ctx, operation.ID)
}

// processStepListUsers lists all existing users
func (p *OperationProcessor) processStepListUsers(ctx context.Context, operation *database.Operation) error {
	log.Printf("Operation %s: Listing existing users", operation.ID)

	// List all users (assuming reasonable number, pagination would be needed for large datasets)
	resp, err := p.userClient.ListUsers(ctx, &userpb.ListUsersRequest{
		Page:  1,
		Limit: 100, // High limit to get all users
	})
	if err != nil {
		// Update operation state to failed
		p.dbClient.UpdateOperationState(ctx, operation.ID, database.StateFailed)
		return fmt.Errorf("failed to list users: %v", err)
	}

	log.Printf("Found %d existing users", len(resp.Users))

	// Store user IDs for next step (in a real implementation, you might store this in operation data)
	// For now, we'll just proceed to deletion

	// Update to next step
	if err := p.dbClient.UpdateOperationStep(ctx, operation.ID, int(StepDeleteUsers), database.StateRunning); err != nil {
		return fmt.Errorf("failed to update operation step: %v", err)
	}

	// Continue processing
	return p.ProcessOperation(ctx, operation.ID)
}

// processStepDeleteUsers deletes all existing users
func (p *OperationProcessor) processStepDeleteUsers(ctx context.Context, operation *database.Operation) error {
	log.Printf("Operation %s: Deleting existing users", operation.ID)

	// List users to get IDs for deletion
	resp, err := p.userClient.ListUsers(ctx, &userpb.ListUsersRequest{
		Page:  1,
		Limit: 100,
	})
	if err != nil {
		p.dbClient.UpdateOperationState(ctx, operation.ID, database.StateFailed)
		return fmt.Errorf("failed to list users for deletion: %v", err)
	}

	// Delete each user
	for _, user := range resp.Users {
		_, err := p.userClient.DeleteUser(ctx, &userpb.DeleteUserRequest{
			Id: user.Id,
		})
		if err != nil {
			// Check if it's a not found error (user might have been deleted already)
			if status.Code(err) != codes.NotFound {
				p.dbClient.UpdateOperationState(ctx, operation.ID, database.StateFailed)
				return fmt.Errorf("failed to delete user %s: %v", user.Id, err)
			}
		}
	}

	log.Printf("Deleted %d users", len(resp.Users))

	// Update to next step
	if err := p.dbClient.UpdateOperationStep(ctx, operation.ID, int(StepCreateUsers), database.StateRunning); err != nil {
		return fmt.Errorf("failed to update operation step: %v", err)
	}

	// Continue processing
	return p.ProcessOperation(ctx, operation.ID)
}

// processStepCreateUsers creates 5 new users
func (p *OperationProcessor) processStepCreateUsers(ctx context.Context, operation *database.Operation) error {
	log.Printf("Operation %s: Creating new users", operation.ID)

	// Create 5 new users
	users := []struct {
		name  string
		email string
		age   int32
	}{
		{"User One", "user1@example.com", 25},
		{"User Two", "user2@example.com", 30},
		{"User Three", "user3@example.com", 35},
		{"User Four", "user4@example.com", 28},
		{"User Five", "user5@example.com", 32},
	}

	for i, userData := range users {
		_, err := p.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
			Name:  userData.name,
			Email: userData.email,
			Age:   userData.age,
		})
		if err != nil {
			p.dbClient.UpdateOperationState(ctx, operation.ID, database.StateFailed)
			return fmt.Errorf("failed to create user %d: %v", i+1, err)
		}
	}

	log.Printf("Created %d new users", len(users))

	// Mark operation as completed
	if err := p.dbClient.UpdateOperationStep(ctx, operation.ID, int(StepCompleted), database.StateCompleted); err != nil {
		return fmt.Errorf("failed to mark operation as completed: %v", err)
	}

	log.Printf("Operation %s completed successfully", operation.ID)
	return nil
}

// StartBackgroundProcessor starts the background worker that processes operations
func (p *OperationProcessor) StartBackgroundProcessor(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(5 * time.Second) // Check for new operations every 5 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Background processor stopped")
				return
			case <-ticker.C:
				p.processPendingOperations(ctx)
			}
		}
	}()
}

// processPendingOperations finds and processes pending operations
func (p *OperationProcessor) processPendingOperations(ctx context.Context) {
	// In a real implementation, you would query for pending operations
	// For this example, we'll rely on the operation being passed directly
	// This method would typically scan the database for operations in PENDING or RUNNING state
	// and resume their processing
}

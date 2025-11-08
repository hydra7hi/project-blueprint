//go:build system
// +build system

package main

import (
	"context"
	"log"
	"testing"
	"time"

	opCl "grpc-services/operation/client"
	opPb "grpc-services/operation/proto"
	userCl "grpc-services/user/client"
	userPb "grpc-services/user/proto"

	_ "github.com/lib/pq"
)

// TestSystem_CompleteWorkflow
// End-to-end test simulating real user scenario:
// 1. User service basic operations
// 2. Operation service processing
// 3. Verify data consistency across services
// 4. Clean up entire system state
func TestSystem_CompleteWorkflow(t *testing.T) {
	ctx := context.Background()

	// Initialize clients for both services
	userClient, err := userCl.NewGRPCClient(ctx, "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create user gRPC client: %v", err)
	}
	defer userClient.Close()

	opClient, err := opCl.NewGRPCClient(ctx, "localhost:50052")
	if err != nil {
		log.Fatalf("Failed to create operation gRPC client: %v", err)
	}
	defer opClient.Close()

	// Phase 1: Initial User Management
	t.Log("Phase 1: Testing User Service Operations")

	// Verify empty initial state
	initialUsers, err := userClient.ListUsers(ctx, &userPb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("Initial ListUsers failed: %s", err)
	}
	if len(initialUsers.GetUsers()) != 0 {
		t.Fatalf("Expected empty database, found %d users", len(initialUsers.GetUsers()))
	}

	// Create test users manually
	_, err = userClient.CreateUser(ctx, &userPb.CreateUserRequest{
		Name:  "System Test User 1",
		Email: "system1@example.com",
		Age:   25,
	})
	if err != nil {
		t.Fatalf("CreateUser 1 failed: %s", err)
	}

	_, err = userClient.CreateUser(ctx, &userPb.CreateUserRequest{
		Name:  "System Test User 2",
		Email: "system2@example.com",
		Age:   30,
	})
	if err != nil {
		t.Fatalf("CreateUser 2 failed: %s", err)
	}

	// Verify users were created
	usersBeforeOp, err := userClient.ListUsers(ctx, &userPb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("ListUsers before operation failed: %s", err)
	}
	if len(usersBeforeOp.GetUsers()) != 2 {
		t.Fatalf("Expected 2 users before operation, got %d", len(usersBeforeOp.GetUsers()))
	}

	// Phase 2: Operation Processing
	t.Log("Phase 2: Testing Operation Service")

	// Start background operation
	opResp, err := opClient.StartOperation(ctx, &opPb.StartOperationRequest{
		OperationData: &opPb.OperationData{},
	})
	if err != nil {
		t.Fatalf("StartOperation failed: %s", err)
	}
	operationID := opResp.GetOperationId()

	// Monitor operation progress
	var operationState string
	for i := 0; i < 20; i++ { // 60 second timeout
		checkResp, err := opClient.CheckProcess(ctx, &opPb.CheckProcessRequest{
			OperationId: operationID,
		})
		if err != nil {
			t.Fatalf("CheckProcess failed: %s", err)
		}

		operationState = checkResp.GetState()
		t.Logf("Operation state: %s (attempt %d)", operationState, i+1)

		if operationState == "COMPLETED" {
			break
		}
		if operationState == "FAILED" {
			t.Fatalf("Operation failed: %s", operationState)
		}

		time.Sleep(3 * time.Second)
	}

	if operationState != "COMPLETED" {
		t.Fatalf("Operation did not complete within timeout. Final state: %s", operationState)
	}

	// Phase 3: Verify System State After Operation
	t.Log("Phase 3: Verifying System State")

	// Operation should have deleted existing users and created 5 new ones
	usersAfterOp, err := userClient.ListUsers(ctx, &userPb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("ListUsers after operation failed: %s", err)
	}

	if len(usersAfterOp.GetUsers()) != 5 {
		t.Fatalf("Expected exactly 5 users after operation, got %d", len(usersAfterOp.GetUsers()))
	}

	// Verify all new users have expected structure
	for i, user := range usersAfterOp.GetUsers() {
		if user.GetName() == "" {
			t.Errorf("User %d has empty name", i+1)
		}
		if user.GetEmail() == "" {
			t.Errorf("User %d has empty email", i+1)
		}
		if user.GetAge() <= 0 {
			t.Errorf("User %d has invalid age: %d", i+1, user.GetAge())
		}

		t.Logf("   %s. %s (%s), Age: %d\n",
			user.GetId(), user.GetName(), user.GetEmail(), user.GetAge())
	}

	// Phase 4: Cross-Service Validation
	t.Log("Phase 4: Cross-Service Validation")

	// Verify operation history
	latestOp, err := opClient.CheckProcess(ctx, &opPb.CheckProcessRequest{
		OperationId: "", // Get latest
	})
	if err != nil {
		t.Fatalf("Get latest operation failed: %s", err)
	}

	if latestOp.GetOperationId() != operationID {
		t.Errorf("Latest operation ID mismatch: got %s, want %s",
			latestOp.GetOperationId(), operationID)
	}

	// Phase 5: System Cleanup
	t.Log("Phase 5: System Cleanup")

	// Delete all users to restore initial state
	for _, user := range usersAfterOp.GetUsers() {
		_, err := userClient.DeleteUser(ctx, &userPb.DeleteUserRequest{
			Id: user.GetId(),
		})
		if err != nil {
			t.Fatalf("DeleteUser failed for %s: %s", user.GetId(), err)
		}
	}

	// Final verification - system should be clean
	finalUsers, err := userClient.ListUsers(ctx, &userPb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("Final ListUsers failed: %s", err)
	}

	if len(finalUsers.GetUsers()) != 0 {
		t.Fatalf("System not clean: %d users remaining", len(finalUsers.GetUsers()))
	}

	t.Log("System test completed successfully - all services working together")
}

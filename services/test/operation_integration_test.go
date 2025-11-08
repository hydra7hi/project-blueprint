//go:build integration
// +build integration

package main

import (
	"context"
	"log"
	"testing"
	"time"

	opCl "grpc-services/operation/client"
	"grpc-services/operation/database"
	opPb "grpc-services/operation/proto"
	userCl "grpc-services/user/client"
	userPb "grpc-services/user/proto"

	_ "github.com/lib/pq"
)

// TestOperation_Integration
// Tests all operation service rpcs.
// The Integration test should run on empty Operations DB.
// And finishes with Empty Operations DB when successful.
//
// The flow:
// 1. Start operation
// 2. Check process status
// 3. Verify users were processed
// 4. Clean up
func TestOperation_Integration(t *testing.T) {
	ctx := context.Background()

	// Expected responses
	expectedOperationStateString := database.StateCompleted.String()

	// Create clients
	opClient, err := opCl.NewGRPCClient(ctx, "localhost:50052")
	if err != nil {
		log.Fatalf("Failed to create operation gRPC client: %v", err)
	}
	defer opClient.Close()

	userClient, err := userCl.NewGRPCClient(ctx, "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create user gRPC client: %v", err)
	}
	defer userClient.Close()

	// 1. Start operation
	respStart, err := opClient.StartOperation(ctx, &opPb.StartOperationRequest{
		OperationData: &opPb.OperationData{},
	})
	if err != nil {
		t.Fatalf("StartOperation error: %s", err)
	}
	operationID := respStart.GetOperationId()
	if operationID == "" {
		t.Fatal("StartOperation returned empty operation ID")
	}

	// 2. Check process status (poll until completed)
	var respCheck *opPb.CheckProcessResponse
	for i := 0; i < 10; i++ { // 10 attempts with delay
		respCheck, err = opClient.CheckProcess(ctx, &opPb.CheckProcessRequest{
			OperationId: operationID,
		})
		if err != nil {
			t.Fatalf("CheckProcess error: %s", err)
		}

		if respCheck.GetState() == expectedOperationStateString {
			break
		}
		// Wait before next check
		time.Sleep(3 * time.Second)
	}

	// Verify operation completed
	if respCheck.GetState() != expectedOperationStateString {
		t.Fatalf("Operation did not complete. Final state: %s", respCheck.GetState())
	}

	// 3. Verify users were processed
	respList, err := userClient.ListUsers(ctx, &userPb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("ListUsers error: %s", err)
	}

	// Should have exactly 5 users created by the operation
	if len(respList.GetUsers()) != 5 {
		t.Fatalf("Expected 5 users, got %d", len(respList.GetUsers()))
	}

	// 4. Clean up - delete all users
	for _, user := range respList.GetUsers() {
		_, err := userClient.DeleteUser(ctx, &userPb.DeleteUserRequest{
			Id: user.GetId(),
		})
		if err != nil {
			t.Fatalf("DeleteUser error: %s", err)
		}
	}

	// Verify clean state
	respList, err = userClient.ListUsers(ctx, &userPb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("ListUsers error: %s", err)
	}
	if len(respList.GetUsers()) != 0 {
		t.Fatalf("Expected 0 users after cleanup, got %d", len(respList.GetUsers()))
	}
}

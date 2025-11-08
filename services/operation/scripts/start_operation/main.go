package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc-services/operation/client"
	pb "grpc-services/operation/proto"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	// Create gRPC client
	gClient, err := client.NewGRPCClient(ctx, "localhost:50052")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer gClient.Close()

	// Start a new operation
	resp, err := gClient.Client.StartOperation(ctx, &pb.StartOperationRequest{
		OperationData: &pb.OperationData{},
	})

	if err != nil {
		panic(fmt.Errorf("failed to start operation: %v\n", err))
	}

	operationID := resp.GetOperationId()
	fmt.Printf("Operation started successfully\n")
	fmt.Printf("   Operation ID: %s\n", operationID)
	fmt.Printf("\nYou can check status with: go run scripts/check_operation.go %s\n", operationID)

	// Wait a moment and show initial status
	time.Sleep(2 * time.Second)
	checkOperationStatus(ctx, gClient, operationID)
}

func checkOperationStatus(ctx context.Context, client *client.GRPCClient, operationID string) {
	resp, err := client.Client.CheckProcess(ctx, &pb.CheckProcessRequest{
		OperationId: operationID,
	})

	if err != nil {
		fmt.Printf("Failed to check operation status: %v\n", err)
		return
	}

	fmt.Printf("\nInitial Status:\n")
	fmt.Printf("   Operation ID: %s\n", resp.GetOperationId())
	fmt.Printf("   Current Step: %d/%d\n", resp.GetCurrentStep(), resp.GetTotalSteps())
	fmt.Printf("   State: %s\n", resp.GetState())
	fmt.Printf("   Completed: %t\n", resp.GetCompleted())
}

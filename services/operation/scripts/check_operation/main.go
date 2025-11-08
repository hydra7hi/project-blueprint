package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"grpc-services/operation/client"
	pb "grpc-services/operation/proto"

	_ "github.com/lib/pq"
)

func main() {
	// Check if operation ID is provided
	if len(os.Args) > 2 {
		fmt.Println("Usage: go run scripts/check_operation.go <operation_id>")
		fmt.Println("Example: go run scripts/check_operation.go op-123456789")
		os.Exit(1)
	}

	operationID := ""
	if len(os.Args) > 1 {
		operationID = os.Args[1]
	}

	ctx := context.Background()

	// Create gRPC client
	gClient, err := client.NewGRPCClient(ctx, "localhost:50052")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer gClient.Close()

	// Check operation status in a loop until completed
	fmt.Printf("Checking status for operation: %s\n", operationID)
	fmt.Println("Press Ctrl+C to stop monitoring")
	fmt.Println("----------------------------------------")

	for {
		resp, err := gClient.Client.CheckProcess(ctx, &pb.CheckProcessRequest{
			OperationId: operationID,
		})

		if err != nil {
			fmt.Printf("Error checking operation: %v\n", err)
			break
		}

		fmt.Printf("Time: %s\n", time.Now().Format("15:04:05"))
		fmt.Printf("   Current Step: %d/%d\n", resp.GetCurrentStep(), resp.GetTotalSteps())
		fmt.Printf("   State: %s\n", resp.GetState())
		fmt.Printf("   Completed: %t\n", resp.GetCompleted())
		fmt.Println("----------------------------------------")

		if resp.GetCompleted() {
			fmt.Println("Operation finished!")
			break
		}

		// Wait before next check
		time.Sleep(3 * time.Second)
	}
}

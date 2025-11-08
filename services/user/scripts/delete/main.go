package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"grpc-services/user/client"
	pb "grpc-services/user/proto"
)

func main() {
	// Check for args or default
	if len(os.Args) > 2 {
		fmt.Println("Usage: go run main.go <user_id>")
		os.Exit(1)
	}

	userID := "1"
	if len(os.Args) > 1 {
		userID = os.Args[1]
	}

	ctx := context.Background()

	// Use the new GRPCClient
	gClient, err := client.NewGRPCClient(ctx, "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer gClient.Close()

	// Delete the user
	resp, err := gClient.Client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: userID,
	})

	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

	fmt.Printf("User deleted successfully: %v\n", resp.GetSuccess())
}

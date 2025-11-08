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
	// Check if operation ID is provided
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

	// Now get the user
	resp, err := gClient.Client.GetUser(ctx, &pb.GetUserRequest{
		Id: userID,
	})

	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
		return
	}

	user := resp.GetUser()
	fmt.Printf("User retrieved successfully\n")
	fmt.Printf("   ID: %s\n", user.GetId())
	fmt.Printf("   Name: %s\n", user.GetName())
	fmt.Printf("   Email: %s\n", user.GetEmail())
	fmt.Printf("   Age: %d\n", user.GetAge())
}

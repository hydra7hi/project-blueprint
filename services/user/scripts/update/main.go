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

	// Update the user
	respOriginal, err := gClient.Client.GetUser(ctx, &pb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
		return
	}
	originalUser := respOriginal.GetUser()

	// Update the user
	resp, err := gClient.Client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    userID,
		Name:  "Updated User",
		Email: "updated@example.com",
		Age:   30,
	})

	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
		return
	}

	user := resp.GetUser()
	fmt.Printf("User updated successfully\n")
	fmt.Printf("   Name: %s → %s\n", originalUser.GetName(), user.GetName())
	fmt.Printf("   Email: %s → %s\n", originalUser.GetEmail(), user.GetEmail())
	fmt.Printf("   Age: %d → %d\n", originalUser.GetAge(), user.GetAge())
	fmt.Printf("   Updated At: %s\n", user.GetUpdatedAt())
}

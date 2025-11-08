package main

import (
	"context"
	"fmt"
	"log"

	"grpc-services/user/client"
	pb "grpc-services/user/proto"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	// Use the new GRPCClient
	gClient, err := client.NewGRPCClient(ctx, "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer gClient.Close()

	// Use the client
	resp, err := gClient.Client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   30,
	})

	if err != nil {
		panic(fmt.Errorf("failed to create user: %v\n", err))
	}

	user := resp.GetUser()
	fmt.Printf("User created successfully\n")
	fmt.Printf("   ID: %s\n", user.GetId())
	fmt.Printf("   Name: %s\n", user.GetName())
	fmt.Printf("   Email: %s\n", user.GetEmail())
	fmt.Printf("   Age: %d\n", user.GetAge())
	fmt.Printf("   Created At: %s\n", user.GetCreatedAt())
	fmt.Printf("   Updated At: %s\n", user.GetUpdatedAt())
}

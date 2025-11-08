package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"grpc-services/user/client"
	pb "grpc-services/user/proto"
)

func main() {
	// Check for args or default
	if len(os.Args) > 3 {
		fmt.Println("Usage: go run main.go <page> <limit>")
		os.Exit(1)
	}

	page := 1
	limit := 10
	var err error
	if len(os.Args) > 1 {
		page, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Usage: go run main.go <page> <limit>")
			os.Exit(1)
		}
		limit, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Usage: go run main.go <page> <limit>")
			os.Exit(1)
		}
	}

	ctx := context.Background()

	// Use the new GRPCClient
	gClient, err := client.NewGRPCClient(ctx, "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer gClient.Close()

	// List users
	resp, err := gClient.Client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})

	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
		return
	}

	fmt.Printf("Listed users successfully\n")
	fmt.Printf("   Total: %d\n", resp.GetTotal())
	fmt.Printf("   Page: %d\n", resp.GetPage())
	fmt.Printf("   Limit: %d\n", resp.GetLimit())
	fmt.Printf("   Users:\n")

	for _, user := range resp.GetUsers() {
		fmt.Printf("   %s. %s (%s), Age: %d\n",
			user.GetId(), user.GetName(), user.GetEmail(), user.GetAge())
	}
}

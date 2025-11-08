package main

import (
	"context"
	"log"

	"grpc-services/operation/client"
	opPb "grpc-services/operation/proto"
	userCl "grpc-services/user/client"
	userPb "grpc-services/user/proto"
)

func main() {
	// This ensures cleanup runs before integration tests
	CleanupDatabases()
}

// CleanupDatabases ensures all test databases are empty before tests run
func CleanupDatabases() {
	ctx := context.Background()
	log.Println("Cleaning up test databases...")

	// Clean user service database
	cleanUserDatabase(ctx)

	// Clean operation service database
	cleanOperationDatabase(ctx)

	log.Println("Database cleanup completed")
}

func cleanUserDatabase(ctx context.Context) {
	userClient, err := userCl.NewGRPCClient(ctx, "localhost:50051")
	if err != nil {
		log.Printf("Warning: Failed to connect to user service: %v", err)
		return
	}
	defer userClient.Close()

	// List all users
	resp, err := userClient.ListUsers(ctx, &userPb.ListUsersRequest{Limit: 1000})
	if err != nil {
		log.Printf("Warning: Failed to list users: %v", err)
		return
	}

	// Delete all users
	for _, user := range resp.GetUsers() {
		_, err := userClient.DeleteUser(ctx, &userPb.DeleteUserRequest{Id: user.GetId()})
		if err != nil {
			log.Printf("Warning: Failed to delete user %s: %v", user.GetId(), err)
		}
	}

	log.Printf("Cleaned %d users from database", len(resp.GetUsers()))
}

func cleanOperationDatabase(ctx context.Context) {
	opClient, err := client.NewGRPCClient(ctx, "localhost:50052")
	if err != nil {
		log.Printf("Warning: Failed to connect to operation service: %v", err)
		return
	}
	defer opClient.Close()

	// Get latest operation to check if any exist
	_, err = opClient.CheckProcess(ctx, &opPb.CheckProcessRequest{})
	if err != nil {
		// If no operations exist, this is fine
		log.Println("No operations to clean")
		return
	}

	log.Println("Operations database cleaned")
}

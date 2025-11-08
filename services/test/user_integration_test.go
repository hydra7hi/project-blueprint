//go:build integration
// +build integration

package main

import (
	"context"
	"log"
	"testing"

	"grpc-services/user/client"
	pb "grpc-services/user/proto"

	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/testing/protocmp"
)

// TestUser_Integration
// Tests all user service rpcs.
// The Integration test should run on empty Users DB.
// And finishes with Empty Users DB when successful.
//
// The the follow the flow:
// 1. List Users.
// 2. Create 2 Users.
// 3. Get User by id.
// 4. List Users.
// 5. Delete Users.
func TestUser_Integration(t *testing.T) {
	ctx := context.Background()

	/// Expected
	expectedListUsersResponse1 := &pb.ListUsersResponse{
		Users: []*pb.User{},
		Total: 0,
		Page:  1,
		Limit: 10,
	}
	expectedListUsersResponse2 := &pb.ListUsersResponse{
		Users: []*pb.User{
			{
				Name:  "User Integration test 1",
				Email: "user1@example.com",
				Age:   1,
			},
			{
				Name:  "User Integration test 2",
				Email: "user2@example.com",
				Age:   2,
			},
		},
		Total: 2,
		Page:  1,
		Limit: 10,
	}
	expectedDeleteUserResponse := &pb.DeleteUserResponse{
		Success: true,
	}

	testUser1 := &pb.User{
		Name:  "User Integration test 1",
		Email: "user1@example.com",
		Age:   1,
	}

	// Use the new GRPCClient
	gClient, err := client.NewGRPCClient(ctx, "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer gClient.Close()

	// 1. List Users.
	respListUsers, err := gClient.ListUsers(ctx, &pb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("ListUsers error: %s", err)
	}
	if diff := cmp.Diff(respListUsers, expectedListUsersResponse1, protocmp.Transform()); diff != "" {
		t.Fatalf("ListUsers response: %s", diff)
	}

	// 2. Create 2 Users.
	respCreateUser1, err := gClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "User Integration test 1",
		Email: "user1@example.com",
		Age:   1,
	})
	if err != nil {
		t.Fatalf("CreateUser error: %s", err)
	}
	respCreateUser2, err := gClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "User Integration test 2",
		Email: "user2@example.com",
		Age:   2,
	})
	if err != nil {
		t.Fatalf("CreateUser error: %s", err)
	}

	// 4. List Users.
	respGetUser, err := gClient.GetUser(ctx, &pb.GetUserRequest{
		Id: respCreateUser1.GetUser().GetId(),
	})
	if err != nil {
		t.Fatalf("GetUser error: %s", err)
	}
	if diff := cmp.Diff(respGetUser.GetUser(), testUser1,
		protocmp.Transform(),
		protocmp.IgnoreFields(&pb.User{}, "id", "created_at", "updated_at")); diff != "" {
		t.Fatalf("ListUsers response: %s", diff)
	}

	// 1. List Users.
	respListUsers, err = gClient.ListUsers(ctx, &pb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("ListUsers error: %s", err)
	}
	if diff := cmp.Diff(respListUsers, expectedListUsersResponse2,
		protocmp.Transform(),
		protocmp.IgnoreFields(&pb.User{}, "id", "created_at", "updated_at")); diff != "" {
		t.Fatalf("ListUsers response: %s", diff)
	}
	// 5. Delete Users.
	respDeleteUser, err := gClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: respCreateUser1.GetUser().GetId(),
	})
	if err != nil {
		t.Fatalf("DeleteUser error: %s", err)
	}
	if !respDeleteUser.GetSuccess() {
		t.Fatalf("DeleteUser res: %s", err)
	}
	if diff := cmp.Diff(respDeleteUser, expectedDeleteUserResponse, protocmp.Transform()); diff != "" {
		t.Fatalf("DeleteUser response: %s", diff)
	}
	respDeleteUser, err = gClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: respCreateUser2.GetUser().GetId(),
	})
	if err != nil {
		t.Fatalf("DeleteUser error: %s", err)
	}
	if diff := cmp.Diff(respDeleteUser, expectedDeleteUserResponse, protocmp.Transform()); diff != "" {
		t.Fatalf("DeleteUser response: %s", diff)
	}
}

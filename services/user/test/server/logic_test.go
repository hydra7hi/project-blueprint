package server

import (
	"context"
	"fmt"
	"testing"

	pb "grpc-services/user/proto"
	"grpc-services/user/server"
	dbMock "grpc-services/user/test/database"

	"github.com/stretchr/testify/assert"
)

func TestServer_CreateUser(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	resp, err := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   30,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Test User", resp.User.Name)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, int32(30), resp.User.Age)
}

func TestServer_GetUser(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	// First create a user
	createResp, _ := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   30,
	})

	// Then get the user
	resp, err := srv.GetUser(context.Background(), &pb.GetUserRequest{
		Id: createResp.User.Id,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Test User", resp.User.Name)
}

func TestServer_UpdateUser(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	// Create user
	createResp, _ := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
		Name:  "Original",
		Email: "original@example.com",
		Age:   25,
	})

	// Update user
	resp, err := srv.UpdateUser(context.Background(), &pb.UpdateUserRequest{
		Id:    createResp.User.Id,
		Name:  "Updated",
		Email: "updated@example.com",
		Age:   30,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.User.Name)
	assert.Equal(t, "updated@example.com", resp.User.Email)
	assert.Equal(t, int32(30), resp.User.Age)
}

func TestServer_DeleteUser(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	// Create a user first
	createResp, _ := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
		Name:  "User To Delete",
		Email: "delete@example.com",
		Age:   30,
	})

	// Delete the user
	resp, err := srv.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		Id: createResp.User.Id,
	})

	assert.NoError(t, err)
	assert.True(t, resp.Success)

	// Verify user is gone
	_, err = srv.GetUser(context.Background(), &pb.GetUserRequest{
		Id: createResp.User.Id,
	})
	assert.Error(t, err)
}

func TestServer_DeleteUser_NotFound(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	_, err := srv.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		Id: "999",
	})

	assert.Error(t, err)
}

func TestServer_ListUsers(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	// Create multiple users
	users := []struct {
		name  string
		email string
		age   int32
	}{
		{"User One", "one@example.com", 25},
		{"User Two", "two@example.com", 30},
		{"User Three", "three@example.com", 35},
	}

	for _, u := range users {
		_, err := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
			Name:  u.name,
			Email: u.email,
			Age:   u.age,
		})
		assert.NoError(t, err)
	}

	// List all users
	resp, err := srv.ListUsers(context.Background(), &pb.ListUsersRequest{
		Page:  1,
		Limit: 10,
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(3), resp.Total)
	assert.Len(t, resp.Users, 3)
	assert.Equal(t, "User One", resp.Users[0].Name)
	assert.Equal(t, "User Two", resp.Users[1].Name)
	assert.Equal(t, "User Three", resp.Users[2].Name)
}

func TestServer_ListUsers_Pagination(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	// Create 5 users
	for i := 1; i <= 5; i++ {
		_, err := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   20 + int32(i),
		})
		assert.NoError(t, err)
	}

	// Test first page with limit 2
	resp1, err := srv.ListUsers(context.Background(), &pb.ListUsersRequest{
		Page:  1,
		Limit: 2,
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(5), resp1.Total)
	assert.Equal(t, int32(1), resp1.Page)
	assert.Equal(t, int32(2), resp1.Limit)
	assert.Len(t, resp1.Users, 2)
	assert.Equal(t, "User 1", resp1.Users[0].Name)
	assert.Equal(t, "User 2", resp1.Users[1].Name)

	// Test second page with limit 2
	resp2, err := srv.ListUsers(context.Background(), &pb.ListUsersRequest{
		Page:  2,
		Limit: 2,
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(2), resp2.Page)
	assert.Len(t, resp2.Users, 2)
	assert.Equal(t, "User 3", resp2.Users[0].Name)
	assert.Equal(t, "User 4", resp2.Users[1].Name)
}

func TestServer_ListUsers_Empty(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	resp, err := srv.ListUsers(context.Background(), &pb.ListUsersRequest{
		Page:  1,
		Limit: 10,
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Total)
	assert.Len(t, resp.Users, 0)
}

func TestServer_ListUsers_DefaultPagination(t *testing.T) {
	mockDB := dbMock.NewMockClient(nil, nil, nil)
	srv := &server.Server{DB: mockDB}

	// Create 15 users to test default limit
	for i := 1; i <= 15; i++ {
		_, err := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   20 + int32(i),
		})
		assert.NoError(t, err)
	}

	// Test with empty request (should use defaults)
	resp, err := srv.ListUsers(context.Background(), &pb.ListUsersRequest{})

	assert.NoError(t, err)
	assert.Equal(t, int32(15), resp.Total)
	assert.Equal(t, int32(1), resp.Page)
	assert.Equal(t, int32(10), resp.Limit) // Default limit
	assert.Len(t, resp.Users, 10)
}

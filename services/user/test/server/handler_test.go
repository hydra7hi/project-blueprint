package server

import (
	"context"
	"errors"
	"testing"

	pb "grpc-services/user/proto"
	"grpc-services/user/server"
	dbMock "grpc-services/user/test/database"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_CreateUser_Handler(t *testing.T) {
	tests := []struct {
		name          string
		givenReq      *pb.CreateUserRequest
		givenDBError  error
		wantErrorCode codes.Code
		wantErrorMsg  string
	}{
		{
			name:     "works - successful creation",
			givenReq: fixtureCreateRequest(),
		},
		{
			name:          "db error - database failure",
			givenReq:      fixtureCreateRequest(),
			givenDBError:  errors.New("database error"),
			wantErrorCode: codes.Internal,
			wantErrorMsg:  "failed to create user",
		},
		{
			name: "validation error - empty name",
			givenReq: fixtureCreateRequest(
				func(req *pb.CreateUserRequest) {
					req.Name = ""
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "name cannot be empty",
		},
		{
			name: "validation error - empty email",
			givenReq: fixtureCreateRequest(
				func(req *pb.CreateUserRequest) {
					req.Email = ""
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "email cannot be empty",
		},
		{
			name: "validation error - invalid email format",
			givenReq: fixtureCreateRequest(
				func(req *pb.CreateUserRequest) {
					req.Email = "invalid-email"
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "invalid email format",
		},
		{
			name: "validation error - zero age",
			givenReq: fixtureCreateRequest(
				func(req *pb.CreateUserRequest) {
					req.Age = 0
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "age must be positive",
		},
		{
			name: "validation error - negative age",
			givenReq: fixtureCreateRequest(
				func(req *pb.CreateUserRequest) {
					req.Age = -5
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "age must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &server.Server{
				DB: dbMock.NewMockClient(
					tt.givenDBError,
					nil,
					nil,
				),
			}

			resp, err := srv.CreateUser(context.Background(), tt.givenReq)

			if tt.wantErrorMsg != "" {
				assert.Error(t, err)
				assert.Nil(t, resp)

				grpcStatus, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrorCode, grpcStatus.Code())
				assert.Contains(t, grpcStatus.Message(), tt.wantErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.givenReq.Name, resp.User.Name)
				assert.Equal(t, tt.givenReq.Email, resp.User.Email)
				assert.Equal(t, tt.givenReq.Age, resp.User.Age)
			}
		})
	}
}

func TestServer_GetUser_Handler(t *testing.T) {
	tests := []struct {
		name          string
		givenReq      *pb.GetUserRequest
		setupMock     func(ctx context.Context, m *dbMock.MockClient)
		wantErrorCode codes.Code
		wantErrorMsg  string
	}{
		{
			name:     "works - successful retrieval",
			givenReq: fixtureGetRequest("1"),
			setupMock: func(ctx context.Context, m *dbMock.MockClient) {
				m.CreateUser(ctx, "Test User", "test@example.com", 30)
			},
		},
		{
			name:          "validation error - empty id",
			givenReq:      fixtureGetRequest(""),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "user ID cannot be empty",
		},
		{
			name:          "db error - user not found",
			givenReq:      fixtureGetRequest("999"),
			wantErrorCode: codes.NotFound,
			wantErrorMsg:  "user not found",
		},
		{
			name:          "validation error - invalid id format",
			givenReq:      fixtureGetRequest("invalid"),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockDB := dbMock.NewMockClient(
				nil,
				nil, nil)
			if tt.setupMock != nil {
				tt.setupMock(ctx, mockDB)
			}
			srv := &server.Server{DB: mockDB}

			resp, err := srv.GetUser(context.Background(), tt.givenReq)

			if tt.wantErrorMsg != "" {
				assert.Error(t, err)
				assert.Nil(t, resp)
				grpcStatus, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrorCode, grpcStatus.Code())
				assert.Contains(t, grpcStatus.Message(), tt.wantErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.User.Id)
			}
		})
	}
}

func TestServer_UpdateUser_Handler(t *testing.T) {
	tests := []struct {
		name          string
		givenReq      *pb.UpdateUserRequest
		setupMock     func(ctx context.Context, m *dbMock.MockClient) string
		wantErrorCode codes.Code
		wantErrorMsg  string
	}{
		{
			name:     "works - successful update",
			givenReq: fixtureUpdateRequest("1"),
			setupMock: func(ctx context.Context, m *dbMock.MockClient) string {
				user, _ := m.CreateUser(ctx, "Original", "original@example.com", 25)
				return user.ToProto().Id
			},
		},
		{
			name:          "validation error - empty id",
			givenReq:      fixtureUpdateRequest(""),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "user ID cannot be empty",
		},
		{
			name:          "validation error - invalid id",
			givenReq:      fixtureUpdateRequest("invalid"),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "invalid user ID",
		},
		{
			name: "validation error - empty name",
			givenReq: fixtureUpdateRequest("1",
				func(req *pb.UpdateUserRequest) {
					req.Name = ""
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "name cannot be empty",
		},
		{
			name: "validation error - empty email",
			givenReq: fixtureUpdateRequest("1",
				func(req *pb.UpdateUserRequest) {
					req.Email = ""
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "email cannot be empty",
		},
		{
			name: "validation error - invalid email format",
			givenReq: fixtureUpdateRequest("1",
				func(req *pb.UpdateUserRequest) {
					req.Email = "invalid-email"
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "invalid email format",
		},
		{
			name: "validation error - zero age",
			givenReq: fixtureUpdateRequest("1",
				func(req *pb.UpdateUserRequest) {
					req.Age = 0
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "age must be positive",
		},
		{
			name:          "db error - user not found",
			givenReq:      fixtureUpdateRequest("999"),
			wantErrorCode: codes.NotFound,
			wantErrorMsg:  "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockDB := dbMock.NewMockClient(
				nil,
				nil, nil)
			var userID string
			if tt.setupMock != nil {
				userID = tt.setupMock(ctx, mockDB)
				if tt.givenReq.Id == "1" {
					tt.givenReq.Id = userID
				}
			}
			srv := &server.Server{DB: mockDB}

			resp, err := srv.UpdateUser(context.Background(), tt.givenReq)

			if tt.wantErrorMsg != "" {
				assert.Error(t, err)
				assert.Nil(t, resp)
				grpcStatus, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrorCode, grpcStatus.Code())
				assert.Contains(t, grpcStatus.Message(), tt.wantErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.givenReq.Name, resp.User.Name)
				assert.Equal(t, tt.givenReq.Email, resp.User.Email)
				assert.Equal(t, tt.givenReq.Age, resp.User.Age)
			}
		})
	}
}

func TestServer_DeleteUser_Handler(t *testing.T) {
	tests := []struct {
		name          string
		givenReq      *pb.DeleteUserRequest
		setupMock     func(ctx context.Context, m *dbMock.MockClient) string
		wantErrorCode codes.Code
		wantErrorMsg  string
	}{
		{
			name:     "works - successful deletion",
			givenReq: fixtureDeleteRequest("1"),
			setupMock: func(ctx context.Context, m *dbMock.MockClient) string {
				user, _ := m.CreateUser(ctx, "To Delete", "delete@example.com", 25)
				return user.ToProto().Id
			},
		},
		{
			name:          "validation error - empty id",
			givenReq:      fixtureDeleteRequest(""),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "user ID cannot be empty",
		},
		{
			name:          "db error - user not found",
			givenReq:      fixtureDeleteRequest("999"),
			wantErrorCode: codes.NotFound,
			wantErrorMsg:  "user not found",
		},
		{
			name:          "validation error - invalid id format",
			givenReq:      fixtureDeleteRequest("invalid"),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockDB := dbMock.NewMockClient(
				nil,
				nil, nil)
			var userID string
			if tt.setupMock != nil {
				userID = tt.setupMock(ctx, mockDB)
				if tt.givenReq.Id == "1" {
					tt.givenReq.Id = userID
				}
			}
			srv := &server.Server{DB: mockDB}

			resp, err := srv.DeleteUser(context.Background(), tt.givenReq)

			if tt.wantErrorMsg != "" {
				assert.Error(t, err)
				assert.Nil(t, resp)
				grpcStatus, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrorCode, grpcStatus.Code())
				assert.Contains(t, grpcStatus.Message(), tt.wantErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
			}
		})
	}
}

func TestServer_ListUsers_Handler(t *testing.T) {
	tests := []struct {
		name              string
		givenReq          *pb.ListUsersRequest
		setupMock         func(ctx context.Context, m *dbMock.MockClient)
		givenDBErrorList  error
		givenDBErrorCount error
		wantErrorCode     codes.Code
		wantErrorMsg      string
		wantUsers         []*pb.User
	}{
		{
			name:     "works - successful list",
			givenReq: fixtureListRequest(),
			setupMock: func(ctx context.Context, m *dbMock.MockClient) {
				m.CreateUser(ctx, "User", "user@example.com", 25)
				m.CreateUser(ctx, "User", "user@example.com", 30)
			},
			wantUsers: []*pb.User{
				fixtureUser(
					"1",
					func(user *pb.User) {
						user.Age = 25
					}),
				fixtureUser(
					"2",
					func(user *pb.User) {
						user.Age = 30
					}),
			},
		},
		{
			name: "works - with pagination",
			givenReq: fixtureListRequest(
				func(req *pb.ListUsersRequest) {
					req.Page = 1
					req.Limit = 1
				}),
			setupMock: func(ctx context.Context, m *dbMock.MockClient) {
				m.CreateUser(ctx, "User", "user@example.com", 25)
				m.CreateUser(ctx, "User 2", "user2@example.com", 30)
			},
			wantUsers: []*pb.User{
				fixtureUser(
					"1",
					func(user *pb.User) {
						user.Age = 25
					}),
			},
		},
		{
			name: "works - with pagination oversized",
			givenReq: fixtureListRequest(
				func(req *pb.ListUsersRequest) {
					req.Page = 1
					req.Limit = 101
				}),
			setupMock: func(ctx context.Context, m *dbMock.MockClient) {
				m.CreateUser(ctx, "User", "user@example.com", 25)
				m.CreateUser(ctx, "User", "user@example.com", 30)
			},
			wantUsers: []*pb.User{
				fixtureUser(
					"1",
					func(user *pb.User) {
						user.Age = 25
					}),
				fixtureUser(
					"2",
					func(user *pb.User) {
						user.Age = 30
					}),
			},
		},
		{
			name: "validation error - negative page",
			givenReq: fixtureListRequest(
				func(req *pb.ListUsersRequest) {
					req.Page = -1
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "page cannot be negative",
		},
		{
			name: "validation error - negative limit",
			givenReq: fixtureListRequest(
				func(req *pb.ListUsersRequest) {
					req.Limit = -1
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "limit cannot be negative",
		},
		{
			name:             "DB error - list",
			givenReq:         fixtureListRequest(),
			givenDBErrorList: errors.New("database error"),
			wantErrorCode:    codes.Internal,
			wantErrorMsg:     "failed to list users",
		},
		{
			name:              "DB error - Count",
			givenReq:          fixtureListRequest(),
			givenDBErrorCount: errors.New("database error"),
			wantErrorCode:     codes.Internal,
			wantErrorMsg:      "failed to count users",
		},
		{
			name: "validation error - negative limit",
			givenReq: fixtureListRequest(
				func(req *pb.ListUsersRequest) {
					req.Limit = -1
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "limit cannot be negative",
		},
		{
			name:      "works - empty list",
			givenReq:  fixtureListRequest(),
			wantUsers: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockDB := dbMock.NewMockClient(
				nil,
				tt.givenDBErrorList,
				tt.givenDBErrorCount)
			if tt.setupMock != nil {
				tt.setupMock(ctx, mockDB)
			}
			srv := &server.Server{DB: mockDB}

			resp, err := srv.ListUsers(context.Background(), tt.givenReq)

			if tt.wantErrorMsg != "" {
				assert.Error(t, err)
				assert.Nil(t, resp)
				grpcStatus, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrorCode, grpcStatus.Code())
				assert.Contains(t, grpcStatus.Message(), tt.wantErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantUsers, resp.Users)
			}
		})
	}
}

// Fixture functions
// Useful to reuse values that are common in test with minor modicaitions.
// mods functions, allow for changing single fields for specific cases.
func fixtureCreateRequest(mods ...func(*pb.CreateUserRequest)) *pb.CreateUserRequest {
	val := &pb.CreateUserRequest{
		Name:  "Valid User",
		Email: "valid@example.com",
		Age:   30,
	}

	for _, mod := range mods {
		mod(val)
	}
	return val
}

func fixtureGetRequest(id string) *pb.GetUserRequest {
	return &pb.GetUserRequest{Id: id}
}

func fixtureUser(id string, mods ...func(*pb.User)) *pb.User {
	val := &pb.User{
		Id:        id,
		Name:      "User",
		Email:     "user@example.com",
		Age:       35,
		CreatedAt: "0001-01-01T00:00:00Z",
		UpdatedAt: "0001-01-01T00:00:00Z",
	}

	for _, mod := range mods {
		mod(val)
	}
	return val
}

func fixtureUpdateRequest(id string, mods ...func(*pb.UpdateUserRequest)) *pb.UpdateUserRequest {
	val := &pb.UpdateUserRequest{
		Id:    id,
		Name:  "Updated User",
		Email: "updated@example.com",
		Age:   35,
	}

	for _, mod := range mods {
		mod(val)
	}
	return val
}

func fixtureDeleteRequest(id string) *pb.DeleteUserRequest {
	return &pb.DeleteUserRequest{Id: id}
}

func fixtureListRequest(mods ...func(*pb.ListUsersRequest)) *pb.ListUsersRequest {
	val := &pb.ListUsersRequest{
		Page:  1,
		Limit: 10,
	}

	for _, mod := range mods {
		mod(val)
	}
	return val
}

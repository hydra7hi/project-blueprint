package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"grpc-services/operation/database"
	pb "grpc-services/operation/proto"
	"grpc-services/operation/server"
	dbMock "grpc-services/operation/test/database"
	userPb "grpc-services/user/proto"
	userMock "grpc-services/user/test/client"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_StartOperation_Handler(t *testing.T) {
	tests := []struct {
		name          string
		givenReq      *pb.StartOperationRequest
		givenDBError  error
		wantErrorCode codes.Code
		wantErrorMsg  string
	}{
		{
			name:     "works - successful operation start",
			givenReq: fixtureStartRequest(),
		},
		{
			name:          "db error - database failure",
			givenReq:      fixtureStartRequest(),
			givenDBError:  errors.New("database error"),
			wantErrorCode: codes.Internal,
			wantErrorMsg:  "failed to create operation",
		},
		{
			name: "validation error - nil operation data",
			givenReq: fixtureStartRequest(
				func(req *pb.StartOperationRequest) {
					req.OperationData = nil
				}),
			wantErrorCode: codes.InvalidArgument,
			wantErrorMsg:  "operation data cannot be empty",
		},
		{
			name: "works - empty operation data",
			givenReq: fixtureStartRequest(
				func(req *pb.StartOperationRequest) {
					req.OperationData = &pb.OperationData{}
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dbClient := dbMock.NewMockClient(
				tt.givenDBError,
				nil,
				nil,
			)

			userClient := &userMock.MockGRPCClient{
				ListUsersResponse:  &userPb.ListUsersResponse{},
				CreateUserResponse: &userPb.UserResponse{},
				UpdateUserResponse: &userPb.UserResponse{},
				GetUserResponse:    &userPb.UserResponse{},
				DeleteUserResponse: &userPb.DeleteUserResponse{},
			}

			srv := &server.Server{
				DB:        dbClient,
				Processor: server.NewTestOperationProcessor(dbClient, userClient),
			}

			resp, err := srv.StartOperation(context.Background(), tt.givenReq)

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
				assert.NotEmpty(t, resp.OperationId)
			}
		})
	}
}

func TestServer_CheckProcess_Handler(t *testing.T) {
	tests := []struct {
		name          string
		givenReq      *pb.CheckProcessRequest
		setupMock     func(ctx context.Context, m *dbMock.MockClient)
		givenDBError  error
		wantErrorCode codes.Code
		wantErrorMsg  string
		wantState     database.OperationState
	}{
		{
			name:     "works - successful check with operation ID",
			givenReq: fixtureCheckRequest("op-1"),
			setupMock: func(ctx context.Context, m *dbMock.MockClient) {
				m.CreateOperation(ctx, fixtureOperation())
			},
			wantState: database.StatePending,
		},
		{
			name:     "works - gets latest when empty operation ID",
			givenReq: fixtureCheckRequest(""),
			setupMock: func(ctx context.Context, m *dbMock.MockClient) {
				m.CreateOperation(ctx, fixtureOperation())
				m.CreateOperation(ctx, fixtureOperation(func(o *database.Operation) {
					o.CreatedAt = time.Now()
					o.State = database.StateCompleted
				}))
			},
			wantState: database.StateCompleted,
		},
		{
			name:          "validation error - operation not found",
			givenReq:      fixtureCheckRequest("non-existent"),
			wantErrorCode: codes.NotFound,
			wantErrorMsg:  "operation not found",
		},
		{
			name:          "db error - database failure",
			givenReq:      fixtureCheckRequest("op-1"),
			givenDBError:  errors.New("database error"),
			wantErrorCode: codes.NotFound,
			wantErrorMsg:  "operation not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockDB := dbMock.NewMockClient(
				nil,
				tt.givenDBError,
				nil)
			if tt.setupMock != nil {
				tt.setupMock(ctx, mockDB)
			}
			srv := &server.Server{DB: mockDB}

			resp, err := srv.CheckProcess(context.Background(), tt.givenReq)

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
				assert.Equal(t, tt.wantState.String(), resp.GetState())
			}
		})
	}
}

// Fixture functions
func fixtureStartRequest(mods ...func(*pb.StartOperationRequest)) *pb.StartOperationRequest {
	val := &pb.StartOperationRequest{
		OperationData: &pb.OperationData{},
	}

	for _, mod := range mods {
		mod(val)
	}
	return val
}

func fixtureCheckRequest(operationID string) *pb.CheckProcessRequest {
	return &pb.CheckProcessRequest{
		OperationId: operationID,
	}
}

func fixtureOperation(mods ...func(*database.Operation)) *database.Operation {
	val := &database.Operation{
		ID:        "op-1",
		StepID:    int(server.StepInitial),
		State:     database.StatePending,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	for _, mod := range mods {
		mod(val)
	}
	return val
}

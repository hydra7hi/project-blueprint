package server

import (
	"context"
	"testing"

	"grpc-services/operation/database"
	"grpc-services/operation/server"
	dbMock "grpc-services/operation/test/database"
	userPb "grpc-services/user/proto"
	userMock "grpc-services/user/test/client"

	"github.com/stretchr/testify/assert"
)

func TestOperationProcessor_ProcessOperation(t *testing.T) {
	tests := []struct {
		name           string
		operationID    string
		setupMock      func(ctx context.Context, m *dbMock.MockClient, u *userMock.MockGRPCClient)
		givenDBError   error
		wantError      bool
		wantErrorMsg   string
		wantStepChange bool
	}{
		{
			name:        "works",
			operationID: "op-1",
			setupMock: func(ctx context.Context, m *dbMock.MockClient, u *userMock.MockGRPCClient) {
				m.CreateOperation(ctx, fixtureOperation())
			},
			wantStepChange: true,
		},
		{
			name:        "error - unknown step",
			operationID: "op-5",
			setupMock: func(ctx context.Context, m *dbMock.MockClient, u *userMock.MockGRPCClient) {
				m.CreateOperation(ctx, &database.Operation{
					ID:     "op-5",
					StepID: 99, // Unknown step
					State:  database.StateRunning,
				})
			},
			wantError:    true,
			wantErrorMsg: "unknown step ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			dbClient := dbMock.NewMockClient(nil, tt.givenDBError, nil)
			userClient := &userMock.MockGRPCClient{
				ListUsersResponse:  &userPb.ListUsersResponse{},
				CreateUserResponse: &userPb.UserResponse{},
				UpdateUserResponse: &userPb.UserResponse{},
				GetUserResponse:    &userPb.UserResponse{},
				DeleteUserResponse: &userPb.DeleteUserResponse{},
			}
			if tt.setupMock != nil {
				tt.setupMock(ctx, dbClient, userClient)
			}

			processor := server.NewTestOperationProcessor(dbClient, userClient)

			err := processor.ProcessOperation(ctx, tt.operationID)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrorMsg)
			} else {
				assert.NoError(t, err)

				if tt.wantStepChange {
					// Verify operation was updated
					op, _ := dbClient.GetOperation(ctx, tt.operationID)
					assert.NotNil(t, op)
					// Step should have progressed from initial state
					if op.StepID == 0 { // Was StepInitial
						assert.Equal(t, 1, op.StepID) // Should be StepListUsers now
					}
				}
			}
		})
	}
}

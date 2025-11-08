package database

import (
	"context"
	"fmt"
	"grpc-services/user/database"
)

// Error Client
type MockClient struct {
	givenCreateError error
	givenListError   error
	givenCountError  error

	Users  map[int]*database.UserRow
	NextID int
}

func NewMockClient(givenCreateError error, givenListError error, givenCountError error) *MockClient {
	return &MockClient{
		givenCreateError: givenCreateError,
		givenListError:   givenListError,
		givenCountError:  givenCountError,
		Users:            make(map[int]*database.UserRow),
		NextID:           1,
	}
}

func (m *MockClient) CreateUser(ctx context.Context, name, email string, age int32) (*database.UserRow, error) {
	if m.givenCreateError != nil {
		return nil, m.givenCreateError
	}

	user := &database.UserRow{
		ID:    m.NextID,
		Name:  name,
		Email: email,
		Age:   age,
	}
	m.Users[m.NextID] = user
	m.NextID++
	return user, nil
}

func (m *MockClient) GetUser(ctx context.Context, id int) (*database.UserRow, error) {
	user, exists := m.Users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MockClient) UpdateUser(ctx context.Context, id int, name, email string, age int32) (*database.UserRow, error) {

	user, exists := m.Users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	user.Name = name
	user.Email = email
	user.Age = age
	return user, nil
}

func (m *MockClient) DeleteUser(ctx context.Context, id int) error {
	if _, exists := m.Users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(m.Users, id)
	return nil
}

func (m *MockClient) ListUsers(ctx context.Context, limit, offset int) ([]*database.UserRow, error) {
	if m.givenListError != nil {
		return nil, m.givenListError
	}

	var users []*database.UserRow
	count := 0
	for i := 1; i < m.NextID; i++ {
		if user, exists := m.Users[i]; exists {
			if count >= offset && len(users) < limit {
				users = append(users, user)
			}
			count++
		}
	}
	return users, nil
}

func (m *MockClient) CountUsers(ctx context.Context) (int, error) {
	if m.givenCountError != nil {
		return 0, m.givenCountError
	}
	return len(m.Users), nil
}

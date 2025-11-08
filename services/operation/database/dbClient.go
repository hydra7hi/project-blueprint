package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"grpc-services/operation/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// SQLClient
// Implements DBClientInterface
type SQLClient struct {
	DB *sql.DB
}

// NewPostgresClient
// Creates the connection to a postgress DB.
// Expects config values to be already checked to not be empty.
//
// Returns:
//   - *SQLClient
//
// Error:
//   - Failed to open sql connection.
//   - Failed to ping database
func NewPostgresClient(cfg *config.Config) (*SQLClient, error) {
	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return &SQLClient{DB: db}, nil
}

// CreateOperation creates a new operation in the database
func (c *SQLClient) CreateOperation(ctx context.Context, op *Operation) error {
	query := `
		INSERT INTO operations 
		(id, marshalled_request, step_id, state, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	_, err := c.DB.ExecContext(ctx, query,
		op.ID,
		op.MarshalledRequest,
		op.StepID,
		op.State.String(),
		now,
		now)
	return err
}

// GetOperation retrieves an operation by ID
func (c *SQLClient) GetOperation(ctx context.Context, id string) (*Operation, error) {
	var op Operation
	var stateStr string
	var createdAt, updatedAt time.Time

	query := `
		SELECT id, marshalled_request, step_id, state, created_at, updated_at 
		FROM operations 
		WHERE id = $1`

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&op.ID,
		&op.MarshalledRequest,
		&op.StepID,
		&stateStr,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = op.State.parse(stateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid state in database: %s", stateStr)
	}

	op.CreatedAt = createdAt
	op.UpdatedAt = updatedAt

	return &op, nil
}

// UpdateOperation updates an existing operation
func (c *SQLClient) UpdateOperation(ctx context.Context, op *Operation) error {
	query := `
		UPDATE operations 
		SET marshalled_request = $1, step_id = $2, state = $3, updated_at = $4 
		WHERE id = $5`

	_, err := c.DB.ExecContext(ctx, query,
		op.MarshalledRequest,
		op.StepID,
		op.State.String(),
		time.Now(),
		op.ID)
	return err
}

// UpdateOperationState updates only the state of an operation
func (c *SQLClient) UpdateOperationState(ctx context.Context, id string, state OperationState) error {
	query := `
		UPDATE operations 
		SET state = $1, updated_at = $2 
		WHERE id = $3`

	_, err := c.DB.ExecContext(ctx, query, state.String(), time.Now(), id)
	return err
}

// UpdateOperationStep updates the step and state of an operation
func (c *SQLClient) UpdateOperationStep(ctx context.Context, id string, stepID int, state OperationState) error {
	query := `
		UPDATE operations 
		SET step_id = $1, state = $2, updated_at = $3 
		WHERE id = $4`

	_, err := c.DB.ExecContext(ctx, query, stepID, state.String(), time.Now(), id)
	return err
}

// GetLatestOperation retrieves the most recently created operation
func (c *SQLClient) GetLatestOperation(ctx context.Context) (*Operation, error) {
	var op Operation
	var stateStr string
	var createdAt, updatedAt time.Time

	query := `
		SELECT id, marshalled_request, step_id, state, created_at, updated_at 
		FROM operations 
		ORDER BY created_at DESC 
		LIMIT 1`

	err := c.DB.QueryRowContext(ctx, query).Scan(
		&op.ID,
		&op.MarshalledRequest,
		&op.StepID,
		&stateStr,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = op.State.parse(stateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid state in database: %s", stateStr)
	}

	op.CreatedAt = createdAt
	op.UpdatedAt = updatedAt

	return &op, nil
}

// CreateTables creates the necessary tables for operations
func (c *SQLClient) CreateTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS operations (
		id VARCHAR(36) PRIMARY KEY,
		marshalled_request JSONB NOT NULL,
		step_id INTEGER NOT NULL DEFAULT 0,
		state VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	DROP TRIGGER IF EXISTS update_operations_updated_at ON operations;
	CREATE TRIGGER update_operations_updated_at
		BEFORE UPDATE ON operations
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`

	_, err := c.DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database tables created/verified successfully")
	return nil
}

// Close closes the database connection
func (c *SQLClient) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

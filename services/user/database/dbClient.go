package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"grpc-services/user/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// SQLClient
// Implements SQLClientInterface
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

func (c *SQLClient) CreateUser(ctx context.Context, name, email string, age int32) (*UserRow, error) {
	var user UserRow
	query :=
		`INSERT INTO users 
		(name, email, age) 
		VALUES ($1, $2, $3) 
	    RETURNING id, name, email, age, created_at, updated_at`
	err := c.DB.QueryRowContext(ctx, query, name, email, age).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return &user, err
}

func (c *SQLClient) GetUser(ctx context.Context, id int) (*UserRow, error) {
	var user UserRow
	query :=
		`SELECT id, name, email, age, created_at, updated_at 
		FROM users 
		WHERE id = $1`
	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return &user, err
}

func (c *SQLClient) UpdateUser(ctx context.Context, id int, name, email string, age int32) (*UserRow, error) {
	var user UserRow
	query :=
		`UPDATE users 
		SET name = $1, email = $2, age = $3 WHERE id = $4 
		RETURNING id, name, email, age, created_at, updated_at`
	err := c.DB.QueryRowContext(ctx, query, name, email, age, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return &user, err
}

func (c *SQLClient) DeleteUser(ctx context.Context, id int) error {
	query :=
		`DELETE 
		FROM users 
		WHERE id = $1`
	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (c *SQLClient) ListUsers(ctx context.Context, limit, offset int) ([]*UserRow, error) {
	query :=
		`SELECT id, name, email, age, created_at, updated_at 
	    FROM users 
		ORDER BY id 
		LIMIT $1 OFFSET $2`
	rows, err := c.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*UserRow
	for rows.Next() {
		var user UserRow
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (c *SQLClient) CountUsers(ctx context.Context) (int, error) {
	var count int
	query :=
		`SELECT COUNT(*) 
		FROM users`
	err := c.DB.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (c *SQLClient) CreateTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		age INTEGER NOT NULL,
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

	DROP TRIGGER IF EXISTS update_users_updated_at ON users;
	CREATE TRIGGER update_users_updated_at
		BEFORE UPDATE ON users
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

func (c *SQLClient) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

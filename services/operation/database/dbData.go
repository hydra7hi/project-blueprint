package database

import (
	"encoding/json"
	"fmt"
	"time"
)

// OperationState represents the state of a long-running operation
type OperationState int

const (
	// StatePending indicates the operation is queued but not yet running
	StatePending OperationState = iota
	// StateRunning indicates the operation is currently executing
	StateRunning
	// StateCompleted indicates the operation finished successfully
	StateCompleted
	// StateFailed indicates the operation failed during execution
	StateFailed
	// StateCancelled indicates the operation was cancelled by the user
	StateCancelled
)

// String returns the string representation of the OperationState
func (s OperationState) String() string {
	return [...]string{"PENDING", "RUNNING", "COMPLETED", "FAILED", "CANCELLED"}[s]
}

func (s *OperationState) parse(stateStr string) error {
	switch stateStr {
	case "PENDING":
		*s = StatePending
	case "RUNNING":
		*s = StateRunning
	case "COMPLETED":
		*s = StateCompleted
	case "FAILED":
		*s = StateFailed
	case "CANCELLED":
		*s = StateCancelled
	default:
		return fmt.Errorf("invalid operation state: %s", stateStr)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (s OperationState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *OperationState) UnmarshalJSON(data []byte) error {
	var stateStr string
	if err := json.Unmarshal(data, &stateStr); err != nil {
		return err
	}

	return s.parse(stateStr)
}

// Operation represents a long-running operation stored in the database
type Operation struct {
	ID                string          `json:"id" db:"id"`
	MarshalledRequest json.RawMessage `json:"marshalled_request" db:"marshalled_request"`
	StepID            int             `json:"step_id" db:"step_id"`
	State             OperationState  `json:"state" db:"state"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
}

// TableName returns the name of the table for the Operation model
func (Operation) TableName() string {
	return "operations"
}

// Validate checks if the operation is valid
func (op *Operation) Validate() error {
	if op.ID == "" {
		return fmt.Errorf("operation ID cannot be empty")
	}
	if op.StepID < 0 {
		return fmt.Errorf("step ID cannot be negative")
	}
	return nil
}

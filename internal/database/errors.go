package database

import (
	"fmt"
)

type DatabaseError struct {
	msg string
	Err error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("DatabaseError: %q: %w", e.msg, e.Err)
}

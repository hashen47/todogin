package config

import (
	"fmt"
)

type ConfError struct {
	msg string
	Err error
}

func (e *ConfError) Error() string {
	return fmt.Sprintf("Config Err: %s: %w", e.msg, e.Err)
}

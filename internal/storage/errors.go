package storage

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

func errNotFound(entity string, id int64) error {
	return fmt.Errorf("storage: %s %d: %w", entity, id, ErrNotFound)
}

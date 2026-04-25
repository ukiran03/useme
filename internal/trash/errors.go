package trash

import "fmt"

type TrashError struct {
	Op   string
	Path string
	Err  error
}

func (e *TrashError) Error() string {
	return fmt.Sprintf("ERR: %s '%s': %v", e.Op, e.Path, e.Err)
}

func (e *TrashError) Unwrap() error {
	return e.Err
}

// Usage:
// return &TrashError{Op: "mkdir", Path: "/.Trash-1000", Err: os.ErrPermission}

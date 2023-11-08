package errors

import (
	"errors"
	"fmt"
)

var (
	BadRequestError = errors.New("bad request")
)

type NotFoundError struct {
	obj string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("object %q not found", e.obj)
}

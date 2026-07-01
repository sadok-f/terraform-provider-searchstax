package client

import (
	"fmt"
)

type Error struct {
	context string
	err     error
}

func (c *Error) Error() string {
	return fmt.Sprintf("%s: %s", c.context, c.err.Error())
}

// Unwrap exposes the wrapped error so callers can use errors.As/errors.Is
// to inspect the underlying cause (e.g. an *HTTPStatusError).
func (c *Error) Unwrap() error {
	return c.err
}

type ApiResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

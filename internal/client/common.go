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

type ApiResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

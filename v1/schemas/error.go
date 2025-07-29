package schemas

import "fmt"

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Implement the error interface
func (e CustomError) Error() string {
	return fmt.Sprintf("code %d: %s", e.Code, e.Message)
}

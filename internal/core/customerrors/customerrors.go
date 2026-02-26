package customerrors

import "fmt"

type ErrStatusCode struct {
	StatusCode int
}

func (e ErrStatusCode) Error() string {
	return fmt.Sprintf("status code: %d", e.StatusCode)
}

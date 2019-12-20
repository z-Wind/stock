package api

import "fmt"

type ErrorAuth struct {
	message string
}

func (e *ErrorAuth) Error() string {
	return fmt.Sprintf("%s", e.message)
}

type ErrorGrant struct {
	message string
}

func (e *ErrorGrant) Error() string {
	return fmt.Sprintf("%s", e.message)
}

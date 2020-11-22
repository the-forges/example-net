package cmdparser

import "fmt"

// ErrorType represent the kind of errors the parser typically will have
type ErrorType string

// ErrorType values
const (
	ErrCommandEmpty     ErrorType = "command empty"
	ErrCommandMalformed ErrorType = "command malformed"
	ErrCommandNotFound  ErrorType = "command not found"
)

func errorForErrorType(errType ErrorType) error {
	return fmt.Errorf("%s", errType)
}

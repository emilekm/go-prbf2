package prism

import (
	"bytes"
	"fmt"
	"strconv"
)

type ErrorCode int

const (
	ErrorCodeUnknown         ErrorCode = iota
	ErrorCodeUnauthenticated ErrorCode = iota + 3000
	ErrorCodeIncorectUsernameOrPassword
	ErrorCodeInssuficientPermissions
	ErrorCodeAccountExists
	ErrorCodeOwnAccont
	ErrorCodeSuperUserLastAccount
	ErrorCodeDeletedUser
	ErrorCodeServerVersion
)

var errorSubjects = []Subject{
	SubjectCriticalError,
	SubjectError,
}

type Error struct {
	Code    ErrorCode
	Details string
}

func (e Error) Subject() Subject {
	return SubjectError
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %d, Details: %s", e.Code, e.Details)
}

var (
	ErrUnknown = Error{
		Code: ErrorCodeUnknown,
	}
)

func ErrorMessageToError(msg Message) error {
	parts := bytes.SplitN(msg.Content(), SeparatorField, 2)

	if len(parts) != 2 {
		return fmt.Errorf("invalid error message: %s", msg.Content())
	}

	code, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return fmt.Errorf("invalid error message: %s", msg.Content())
	}

	return Error{
		Code:    ErrorCode(code),
		Details: string(parts[1]),
	}
}

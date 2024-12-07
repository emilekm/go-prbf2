package prism2

import (
	"fmt"
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

func (e Error) Error() string {
	return fmt.Sprintf("Code: %d, Details: %s", e.Code, e.Details)
}

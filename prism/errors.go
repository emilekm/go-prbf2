package prism

import "fmt"

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
	Content string
}

func (e Error) Subject() Subject {
	return SubjectError
}

func (e Error) Error() string {
	return fmt.Sprintf("ErrorCode: %d, Content: %s", e.Code, e.Content)
}

type CriticalError Error

func (e CriticalError) Subject() Subject {
	return SubjectCriticalError
}

package prism

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
	Content string
}

func (e Error) Error() string {
	return fmt.Sprintf("ErrorCode: %d, Content: %s", e.Code, e.Content)
}

func NewError(c ErrorCode, content string) Error {
	return Error{
		Code:    c,
		Content: content,
	}
}

func NewErrorFromMessage(msg Message) error {
	var msgErr Error
	err := UnmarshalInto(msg, &msgErr)
	if err != nil {
		return err
	}

	return msgErr
}

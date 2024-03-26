package prism

import (
	"errors"
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
	Type    Subject
	Code    ErrorCode
	Content string
}

func (e Error) Error() string {
	return fmt.Sprintf("Error: %s, Code: %d, Content: %s", e.Type, e.Code, e.Content)
}

func NewError(t Subject, c ErrorCode, content string) Error {
	return Error{
		Type:    t,
		Code:    c,
		Content: content,
	}
}

func NewErrorFromMessage(msg Message) error {
	if len(msg.Fields) < 2 {
		return errors.New("invalid error message")
	}

	code, err := strconv.Atoi(string(msg.Fields[0]))
	if err != nil {
		code = int(ErrorCodeUnknown)
	}

	return NewError(
		msg.Subject,
		ErrorCode(code),
		string(msg.Fields[1]),
	)
}

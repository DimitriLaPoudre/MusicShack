package model

import (
	"errors"
	"fmt"
)

// -- Message -- //

type Message struct {
	code string
	args []any
}

func (m *Message) String() string { return fmt.Sprintf(string(m.code), m.args...) }

var (
	MessageCodeNotFound           string = "%v not found"
	MessageCodeInvalid            string = "%v is invalid"
	MessageCodeExpired            string = "%v has expired"
	MessageCodeAlreadyExists      string = "%v already exists"
	MessageCodeIncorrect          string = "%v is incorrect"
	MessageCodeEmpty              string = "%v is empty"
	MessageCodePluginConfigFailed string = "plugin configuration failed: %v"
	MessageCodeInternalError      string = "internal error"
)

// -- Error -- //

type Error struct {
	err     error
	message Message
}

func NewError(err error, code string, args ...any) error {
	return Error{
		err: err,
		message: Message{
			code,
			args,
		},
	}
}
func (e Error) Error() string    { return fmt.Sprintf("%s: %v", e.message.String(), e.err) }
func (e *Error) Unwrap() error   { return e.err }
func (e *Error) Message() string { return e.message.String() }

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalid            = errors.New("invalid")
	ErrExpired            = errors.New("expired")
	ErrConflict           = errors.New("conflict")
	ErrIncorrect          = errors.New("incorrect")
	ErrEmpty              = errors.New("empty")
	ErrPluginConfigFailed = errors.New("plugin config failed")
	ErrInternal           = errors.New("internal error")

	ErrRateLimitReached = errors.New("rate limit reached")
	ErrUnknown          = errors.New("unknown error")
)

// -- OutboundError -- //

type OutboundError struct {
	appErr      error
	outboundErr error
}

func NewOutboundError(appErr error, outboundErr error) error {
	return OutboundError{
		appErr:      appErr,
		outboundErr: outboundErr,
	}
}
func (e OutboundError) Error() string  { return e.outboundErr.Error() }
func (e *OutboundError) Unwrap() error { return e.appErr }

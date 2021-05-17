package errors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	// Error is a type that represents
	// the application error with context
	// and an underlying cause.
	Error struct {
		kind    ErrorKind
		context string
		cause   error
	}
	// ErrorKind is a type that represents
	// various kind of error types.
	ErrorKind int
)

// Defines possible error kinds.
const (
	// Internal represents an internal error kind.
	// It is used when an error occurs in an internal service or gateway.
	Internal ErrorKind = http.StatusInternalServerError
	// InvalidInput represents an invalid input error kind.
	// It is used when trying to process bad payload or input data.
	InvalidInput ErrorKind = http.StatusBadRequest
	// NotFound represents a not found error kind.
	// It is used when trying to retrieve a nonexistent entity.
	NotFound ErrorKind = http.StatusNotFound
)

// New creates a new instance of Error.
func New(kind ErrorKind, ctx string) *Error {
	return &Error{
		kind:    kind,
		context: ctx,
		cause:   errors.New(ctx),
	}
}

// Newf creates a new instance of Error
// with a formatted context
func Newf(kind ErrorKind, ctx string, args ...interface{}) *Error {
	return New(kind, fmt.Sprintf(ctx, args...))
}

// Wrap creates a new instance of Error
// with a context and an underlying error.
func Wrap(cause error, kind ErrorKind, ctx string) *Error {
	return &Error{
		kind:    kind,
		context: ctx,
		cause:   errors.WithStack(cause),
	}
}

// Wrapf creates a new instance of Error
// with a formatted context and an underlying cause.
func Wrapf(cause error, kind ErrorKind, ctx string, args ...interface{}) *Error {
	return Wrap(cause, kind, fmt.Sprintf(ctx, args...))
}

// Context returns an internal context.
func (e *Error) Context() string {
	return e.context
}

// Kind returns an internal error kind.
func (e *Error) Kind() ErrorKind {
	return e.kind
}

// Cause returns an underlying error.
func (e *Error) Cause() error {
	return e.cause
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return e.context
}

// New creates a new instance of Error with a context.
func (kind ErrorKind) New(ctx string) error {
	return New(kind, ctx)
}

// Newf creates a new instance of Error with a formatted context.
func (kind ErrorKind) Newf(ctx string, args ...interface{}) error {
	return New(kind, fmt.Sprintf(ctx, args...))
}

// Wrap creates a new instance of Error
// with a context and an underlying cause.
func (kind ErrorKind) Wrap(cause error, ctx string) error {
	return Wrap(cause, kind, ctx)
}

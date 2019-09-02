package errors

import (
	"context"
	"fmt"
	"io"
	"reflect"

	"github.com/felipeweb/devctl/internal/retry"
	"gocloud.dev/gcerrors"
	"golang.org/x/xerrors"
)

// An Error describes a Go CDK error.
type Error struct {
	Code  gcerrors.ErrorCode
	msg   string
	frame xerrors.Frame
	err   error
}

func (e *Error) Error() string {
	return fmt.Sprint(e)
}

func (e *Error) Format(s fmt.State, c rune) {
	xerrors.FormatError(e, s, c)
}

func (e *Error) FormatError(p xerrors.Printer) (next error) {
	if e.msg == "" {
		p.Printf("code=%v", e.Code)
	} else {
		p.Printf("%s (code=%v)", e.msg, e.Code)
	}
	e.frame.Format(p)
	return e.err
}

// Unwrap returns the error underlying the receiver, which may be nil.
func (e *Error) Unwrap() error {
	return e.err
}

// New returns a new error with the given code, underlying error and message. Pass 1
// for the call depth if New is called from the function raising the error; pass 2 if
// it is called from a helper function that was invoked by the original function; and
// so on.
func New(c gcerrors.ErrorCode, err error, callDepth int, msg string) *Error {
	return &Error{
		Code:  c,
		msg:   msg,
		frame: xerrors.Caller(callDepth),
		err:   err,
	}
}

// Newf uses format and args to format a message, then calls New.
func Newf(c gcerrors.ErrorCode, err error, format string, args ...interface{}) *Error {
	return New(c, err, 2, fmt.Sprintf(format, args...))
}

// ErrorAs is a helper for the ErrorAs method of an API's portable type.
// It performs some initial nil checks, and does a single level of unwrapping
// when err is a *gcerr.Error. Then it calls its errorAs argument, which should
// be a driver implementation of ErrorAs.
func ErrorAs(err error, target interface{}, errorAs func(error, interface{}) bool) bool {
	if err == nil {
		return false
	}
	if target == nil {
		panic("ErrorAs target cannot be nil")
	}
	val := reflect.ValueOf(target)
	if val.Type().Kind() != reflect.Ptr || val.IsNil() {
		panic("ErrorAs target must be a non-nil pointer")
	}
	if e, ok := err.(*Error); ok {
		err = e.Unwrap()
	}
	return errorAs(err, target)
}

// DoNotWrap reports whether an error should not be wrapped in the Error
// type from this package.
// It returns true if err is a retry error, a context error, io.EOF, or if it wraps
// one of those.
func DoNotWrap(err error) bool {
	if xerrors.Is(err, io.EOF) {
		return true
	}
	if xerrors.Is(err, context.Canceled) {
		return true
	}
	if xerrors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var r *retry.ContextError
	return xerrors.As(err, &r)
}

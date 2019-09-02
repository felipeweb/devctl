package driver

import (
	"context"

	"gocloud.dev/gcerrors"
)

// Machine interface
type Machine interface {
	// ErrorCode should return a code that describes the error, which was returned by
	// one of the other methods in this interface.
	ErrorCode(error) gcerrors.ErrorCode

	// As converts i to driver-specific types.
	// See https://gocloud.dev/concepts/as/ for background information.
	As(i interface{}) bool

	// ErrorAs allows drivers to expose driver-specific types for returned
	// errors.
	// See https://gocloud.dev/concepts/as/ for background information.
	ErrorAs(error, interface{}) bool

	// Start the machine
	Start(context.Context) error

	// Shutdown the machine
	Shutdown(context.Context) error
}

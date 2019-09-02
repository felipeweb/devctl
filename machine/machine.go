package machine

import (
	"context"
	"fmt"
	"sync"

	"github.com/felipeweb/devctl/internal/errors"
	"github.com/felipeweb/devctl/machine/driver"
	"gocloud.dev/gcerrors"
)

// Machine provides an easy and portable way to interact with vms
type Machine struct {
	m driver.Machine
	// mu protects the closed variable.
	// Read locks are kept to allow holding a read lock for long-running calls,
	// and thereby prevent closing until a call finishes.
	mu  sync.RWMutex
	off bool
}

// As converts i to driver-specific types.
// See https://gocloud.dev/concepts/as/ for background information, the "As"
// examples in this package for examples, and the driver package
// documentation for the specific types supported for that driver.
func (m *Machine) As(i interface{}) bool {
	if i == nil {
		return false
	}
	return m.m.As(i)
}

// ErrorAs converts err to driver-specific types.
// ErrorAs panics if i is nil or not a pointer.
// ErrorAs returns false if err == nil.
// See https://gocloud.dev/concepts/as/ for background information.
func (m *Machine) ErrorAs(err error, i interface{}) bool {
	return errors.ErrorAs(err, i, m.m.ErrorAs)
}

// Shutdown releases any resources used for the machine.
func (m *Machine) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	prev := m.off
	m.off = true
	m.mu.Unlock()
	if prev {
		return errClosed
	}
	return wrapError(m.m, m.m.Shutdown(ctx), "")
}

// Start the machine
func (m *Machine) Start(ctx context.Context) error {
	m.mu.Lock()
	m.off = false
	m.mu.Unlock()
	return wrapError(m.m, m.m.Start(ctx), "")
}

func wrapError(m driver.Machine, err error, name string) error {
	if err == nil {
		return nil
	}
	if errors.DoNotWrap(err) {
		return err
	}
	msg := "vm"
	if name != "" {
		msg += fmt.Sprintf(" (name %q)", name)
	}
	return errors.New(m.ErrorCode(err), err, 2, msg)
}

var errClosed = errors.Newf(gcerrors.FailedPrecondition, nil, "vm: Machine has been shutdown")

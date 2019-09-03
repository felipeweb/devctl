package machine

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/felipeweb/devctl/devenv/machine/driver"
	"github.com/felipeweb/devctl/internal/errors"
	"github.com/felipeweb/devctl/internal/openurl"
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

// MachineURLOpener represents types that can open Machines based on a URL.
// The opener must not modify the URL argument. NewURL must be safe to
// call from multiple goroutines.
//
// This interface is generally implemented by types in driver packages.
type MachineURLOpener interface {
	NewURL(ctx context.Context, u *url.URL) (*Machine, error)
}

// URLMux is a URL opener multiplexer. It matches the scheme of the URLs
// against a set of registered schemes and calls the opener that matches the
// URL's scheme.
// See https://gocloud.dev/concepts/urls/ for more information.
//
// The zero value is a multiplexer with no registered schemes.
type URLMux struct {
	schemes openurl.SchemeMap
}

// MachineSchemes returns a sorted slice of the registered Machine schemes.
func (mux *URLMux) MachineSchemes() []string { return mux.schemes.Schemes() }

// ValidMachineScheme returns true iff scheme has been registered for Machines.
func (mux *URLMux) ValidMachineScheme(scheme string) bool { return mux.schemes.ValidScheme(scheme) }

// RegisterMachine registers the opener with the given scheme. If an opener
// already exists for the scheme, RegisterMachine panics.
func (mux *URLMux) RegisterMachine(scheme string, opener MachineURLOpener) {
	mux.schemes.Register("vm", "machine", scheme, opener)
}

// New calls NewURL with the URL parsed from urlstr.
// New is safe to call from multiple goroutines.
func (mux *URLMux) New(ctx context.Context, urlstr string) (*Machine, error) {
	o, u, err := mux.schemes.FromString("machine", urlstr)
	if err != nil {
		return nil, err
	}
	opener := o.(MachineURLOpener)
	machine, err := opener.NewURL(ctx, u)
	if err != nil {
		return nil, err
	}
	return machine, nil
}

// NewURL dispatches the URL to the opener that is registered with the
// URL's scheme. NewURL is safe to call from multiple goroutines.
func (mux *URLMux) NewURL(ctx context.Context, u *url.URL) (*Machine, error) {
	o, err := mux.schemes.FromURL("machine", u)
	if err != nil {
		return nil, err
	}
	opener := o.(MachineURLOpener)
	machine, err := opener.NewURL(ctx, u)
	if err != nil {
		return nil, err
	}
	return machine, nil
}

var defaultURLMux = new(URLMux)

// DefaultURLMux returns the URLMux used by New.
//
// Driver packages can use this to register their MachineURLOpener on the mux.
func DefaultURLMux() *URLMux {
	return defaultURLMux
}

// New opens the Machine identified by the URL given.
//
// See the URLOpener documentation in driver subpackages for
// details on supported URL formats, and https://gocloud.dev/concepts/urls/
// for more information.
func New(ctx context.Context, urlstr string) (*Machine, error) {
	return defaultURLMux.New(ctx, urlstr)
}

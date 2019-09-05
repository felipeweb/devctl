package context

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/afero"
)

type CommandBuilder func(name string, arg ...string) *exec.Cmd

type Context struct {
	Debug      bool
	verbose    bool
	FileSystem *afero.Afero
	In         io.Reader
	Out        io.Writer
	Err        io.Writer
	NewCommand CommandBuilder
}

func (c *Context) SetVerbose(value bool) {
	c.verbose = value
}

func (c *Context) IsVerbose() bool {
	return c.Debug || c.verbose
}

// CensoredWriter is a writer wrapping the provided io.Writer with logic to censor certain values
type CensoredWriter struct {
	writer          io.Writer
	sensitiveValues []string
}

// NewCensoredWriter returns a new CensoredWriter
func NewCensoredWriter(writer io.Writer) *CensoredWriter {
	return &CensoredWriter{writer: writer, sensitiveValues: []string{}}
}

// SetSensitiveValues sets values needing masking for an CensoredWriter
func (cw *CensoredWriter) SetSensitiveValues(vals []string) {
	cw.sensitiveValues = vals
}

// Write implements io.Writer's Write method, performing necessary auditing while doing so
func (cw *CensoredWriter) Write(b []byte) (int, error) {
	auditedBytes := b
	for _, val := range cw.sensitiveValues {
		auditedBytes = bytes.Replace(auditedBytes, []byte(val), []byte("*******"), -1)
	}

	_, err := cw.writer.Write(auditedBytes)
	return len(b), err
}

func New() *Context {
	// Default to respecting the DEVCTL_DEBUG env variable, the cli will override if --debug is set otherwise
	_, debug := os.LookupEnv("DEVCTL_DEBUG")

	return &Context{
		Debug:      debug,
		FileSystem: &afero.Afero{Fs: afero.NewOsFs()},
		In:         os.Stdin,
		Out:        NewCensoredWriter(os.Stdout),
		Err:        NewCensoredWriter(os.Stderr),
		NewCommand: exec.Command,
	}
}

// SetSensitiveValues sets the sensitive values needing masking on output/err streams
func (c *Context) SetSensitiveValues(vals []string) {
	if len(vals) > 0 {
		out := NewCensoredWriter(c.Out)
		out.SetSensitiveValues(vals)
		c.Out = out

		err := NewCensoredWriter(c.Err)
		err.SetSensitiveValues(vals)
		c.Err = err
	}
}

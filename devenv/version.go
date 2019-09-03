package devenv

import (
	"fmt"

	"github.com/felipeweb/devctl/devenv/printer"
)

var (
	Commit  string
	Version string
)

// VersionOptions represent generic options for use by Porter's list commands
type VersionOptions struct {
	printer.PrintOptions
}

var DefaultVersionFormat = printer.FormatPlaintext

func (o *VersionOptions) Validate() error {
	if o.RawFormat == "" {
		o.RawFormat = string(DefaultVersionFormat)
	}

	err := o.ParseFormat()
	if err != nil {
		return err
	}

	switch o.Format {
	case printer.FormatJson, printer.FormatPlaintext, printer.FormatYaml, printer.FormatTable:
		return nil
	default:
		return fmt.Errorf("unsupported format, %s. Supported formats are: %s, %s, %s, %s", o.Format, printer.FormatJson, printer.FormatPlaintext, printer.FormatYaml, printer.FormatTable)
	}
}

func (d *Devenv) PrintVersion(opts VersionOptions) error {
	type version struct {
		Version string
		Commit  string
	}
	v := version{
		Version,
		Commit,
	}
	switch opts.Format {
	case printer.FormatJson:
		return printer.PrintJson(d.Out, v)
	case printer.FormatYaml:
		return printer.PrintYaml(d.Out, v)
	case printer.FormatPlaintext:
		return printer.PrintPlaintext(d.Out, fmt.Sprintf("Devctl %v (%v)", Version, Commit))
	case printer.FormatTable:
		f := func(row interface{}) []interface{} {
			i, ok := row.(version)
			if !ok {
				return nil
			}
			return []interface{}{i.Version, i.Commit}
		}
		return printer.PrintTable(d.Out, []version{v}, f, "VERSION", "COMMIT")
	default:
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

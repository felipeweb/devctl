package main

import (
	"github.com/felipeweb/devctl/devenv"
	"github.com/spf13/cobra"
)

func buildVersionCommand(d *devenv.Devenv) *cobra.Command {
	opts := devenv.VersionOptions{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return d.PrintVersion(opts)
		},
	}
	cmd.Annotations = map[string]string{
		"group": "meta",
	}
	f := cmd.Flags()
	f.StringVarP(&opts.RawFormat, "output", "o", string(devenv.DefaultVersionFormat),
		"Specify an output format.  Allowed values: yaml, json, table, plaintext")

	return cmd
}

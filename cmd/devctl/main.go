package main

import (
	"os"

	"github.com/felipeweb/devctl/devenv"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRootCommand() *cobra.Command {
	d := devenv.New()
	cmd := &cobra.Command{
		Use:   "devctl",
		Short: "Help developers save money with remote dev environment",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Enable swapping out stdout/stderr for testing
			d.Out = cmd.OutOrStdout()
			d.Err = cmd.OutOrStderr()
		},
		SilenceUsage: true,
	}

	cmd.PersistentFlags().BoolVar(&d.Debug, "debug", false, "Enable debug logging")
	cmd.AddCommand(buildVersionCommand(d))
	for _, alias := range buildAliasCommands(d) {
		cmd.AddCommand(alias)
	}

	help := newHelptextBox()
	usage, _ := help.FindString("usage.tmpl")
	cmd.SetUsageTemplate(usage)
	cobra.AddTemplateFunc("ShouldShowGroupCommands", ShouldShowGroupCommands)
	cobra.AddTemplateFunc("ShouldShowGroupCommand", ShouldShowGroupCommand)
	cobra.AddTemplateFunc("ShouldShowUngroupedCommands", ShouldShowUngroupedCommands)
	cobra.AddTemplateFunc("ShouldShowUngroupedCommand", ShouldShowUngroupedCommand)
	return cmd
}

func newHelptextBox() *packr.Box {
	return packr.New("github.com/felipeweb/devctl/cmd/devctl/helptext", "./helptext")
}

func ShouldShowGroupCommands(cmd *cobra.Command, group string) bool {
	for _, child := range cmd.Commands() {
		if ShouldShowGroupCommand(child, group) {
			return true
		}
	}
	return false
}

func ShouldShowGroupCommand(cmd *cobra.Command, group string) bool {
	return cmd.Annotations["group"] == group
}

func ShouldShowUngroupedCommands(cmd *cobra.Command) bool {
	for _, child := range cmd.Commands() {
		if ShouldShowUngroupedCommand(child) {
			return true
		}
	}
	return false
}

func ShouldShowUngroupedCommand(cmd *cobra.Command) bool {
	if !cmd.IsAvailableCommand() {
		return false
	}

	_, hasGroup := cmd.Annotations["group"]
	return !hasGroup
}

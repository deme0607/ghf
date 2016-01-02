package main

import (
	"flag"
	"fmt"
	"io"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		noEditor bool
		params   string
		force    bool
		issue    string
		base     string
		head     string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&noEditor, "no-editor", false, "Disable opening $EDITOR")
	flags.BoolVar(&noEditor, "n", false, "Disable opening $EDITOR(Short)")
	flags.StringVar(&params, "params", "", "Parameters for template")
	flags.StringVar(&params, "p", "", "Parameters for template(Short)")

	// parameters for hub pull-request command
	flags.BoolVar(&force, "force", false, "")
	flags.BoolVar(&force, "f", false, "(Short)")
	flags.StringVar(&issue, "issue", "", "")
	flags.StringVar(&issue, "i", "", "(Short)")
	flags.StringVar(&base, "base", "", "")
	flags.StringVar(&base, "b", "", "(Short)")
	flags.StringVar(&head, "head", "", "")
	flags.StringVar(&head, "h", "", "(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	_ = noEditor

	_ = params

	_ = force

	_ = issue

	_ = base

	_ = head

	return ExitCodeOK
}

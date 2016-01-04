package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/github/hub/cmd"
	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota

	ConfigFileName string = ".ghf-template"
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

func (cli *CLI) getTemplate() string {
	gitDir, err := git.Dir()
	if err != nil {
		fmt.Fprintf(cli.errStream, err.Error())
		return ""
	}

	config_file := filepath.Join(gitDir, "..", ConfigFileName)

	if _, err := os.Stat(config_file); err != nil {
		current_user, _ := user.Current()
		config_file = filepath.Join(current_user.HomeDir, ConfigFileName)

		if _, err := os.Stat(config_file); err != nil {
			fmt.Fprintf(cli.errStream, err.Error())
			return ""
		}
	}

	text, err := ioutil.ReadFile(config_file)
	if err != nil {
		fmt.Fprintf(cli.errStream, err.Error())
		return ""
	}

	return strings.TrimRight(string(text), "\n")
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		noEditor      bool
		params        string
		template_file string
		force         bool
		issue         string
		base          string
		head          string
		browse        bool

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&noEditor, "no-editor", false, "Disable opening $EDITOR")
	flags.BoolVar(&noEditor, "n", false, "Disable opening $EDITOR(Short)")
	flags.StringVar(&params, "params", "", "Parameters for template")
	flags.StringVar(&params, "p", "", "Parameters for template(Short)")
	flags.StringVar(&template_file, "template", "", "Template filename")
	flags.StringVar(&template_file, "t", "", "Template filename(Short)")

	// parameters for hub pull-request command
	flags.BoolVar(&force, "force", false, "")
	flags.BoolVar(&force, "f", false, "(Short)")
	flags.BoolVar(&browse, "browse", false, "")
	flags.BoolVar(&browse, "o", false, "(Short)")
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

	if template_file == "" {
		template_file = cli.getTemplate()

		if template_file == "" {
			return ExitCodeError
		}
	}

	params_map := make(map[string]string)
	if params != "" {
		for _, p := range strings.Split(params, ",") {
			param := strings.Split(p, ":")
			params_map[param[0]] = param[1]
		}
	}

	tmpl := template.Must(template.ParseFiles(template_file))
	out := new(bytes.Buffer)
	tmpl.Execute(out, params_map)

	content := out.String()

	if !noEditor {
		editor, err := github.NewEditor("PULLREQ", "pull request", content)
		if err != nil {
			return ExitCodeError
		}

		defer editor.DeleteFile()

		editor.CS = "//"

		title, body, err := editor.EditTitleAndBody()
		if err != nil {
			return ExitCodeError
		}

		content = title + "\n\n" + body
	}

	hubCmd := cmd.New("hub pull-request")
	hubCmd.WithArg("-m " + content)
	if force {
		hubCmd.WithArg("-f")
	}
	if browse {
		hubCmd.WithArg("-o")
	}
	if issue != "" {
		hubCmd.WithArg("-i " + issue)
	}
	if base != "" {
		hubCmd.WithArg("-b " + base)
	}
	if head != "" {
		hubCmd.WithArg("-h " + head)
	}

	err := hubCmd.Spawn()
	if err != nil {
		return ExitCodeError
	}

	return ExitCodeOK
}

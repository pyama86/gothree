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
		awsKey string
		awsId  string
		region string
		bucket string
		path   string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&awsKey, "aws-key", "", "")
	flags.StringVar(&awsKey, "a", "", "(Short)")

	flags.StringVar(&awsId, "aws-id", "", "")
	flags.StringVar(&awsId, "a", "", "(Short)")

	flags.StringVar(&region, "region", "", "")
	flags.StringVar(&region, "r", "", "(Short)")

	flags.StringVar(&bucket, "bucket", "", "")
	flags.StringVar(&bucket, "b", "", "(Short)")

	flags.StringVar(&path, "path", "", "")
	flags.StringVar(&path, "p", "", "(Short)")

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

	_ = awsKey

	_ = awsId

	_ = region

	_ = bucket

	_ = path

	return ExitCodeOK
}

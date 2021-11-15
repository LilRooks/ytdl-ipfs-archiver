package app

import (
	"flag"
	"io"
)

// ErrorCode represents the error and code to be returned by the command.
type ErrorCode struct {
	Error error
	Code  int
}

var ytdlBinary string
var configPath string
var tablePath string

// Run is the actual code for the command
func Run(args []string, stdout io.Writer) ErrorCode {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.StringVar(&ytdlBinary, "bin", "", "path to the youtube-dl binary (defaults to one in PATH)")
	flags.StringVar(&configPath, "cfg", "", "path to the configuration file to use")
	flags.StringVar(&tablePath, "tab", "./table.edn", "path to the table file to use")

	if err := flags.Parse(args[1:]); err != nil {
		return ErrorCode{Error: err, Code: -1}
	}
	ytdlOptions := flags.Args()
	return ErrorCode{Error: nil, Code: 0}
}

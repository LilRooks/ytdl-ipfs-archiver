package app

import (
	"flag"
	"fmt"
	"io"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/config"
)

var ytdlPath string
var confPath string
var tablPath string

const (
	errorGeneric = iota
	errorConfig
	errorYTDL
	errorTable
	errorIPFS
)

// Run is the actual code for the command
func Run(args []string, stdout io.Writer) (error, int) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.StringVar(&ytdlPath, "bin", "", "path to the youtube-dl binary (defaults to one in PATH)")
	flags.StringVar(&confPath, "cfg", "", "path to the configuration file to use")
	flags.StringVar(&tablPath, "tab", "./table.sqlite", "path to the table file to use")

	if err := flags.Parse(args[1:]); err != nil {
		return err, errorGeneric
	}

	err, configs := config.Parse(confPath)
	if err != nil {
		return err, errorConfig
	}

	if len(ytdlPath) == 0 {
		ytdlPath = configs.Binary
	}
	if err := ytdl.CheckBinary(ytdlPath); err != nil {
		return err, errorYTDL
	}

	fmt.Fprintf(stdout, "Binary at %s\n", ytdlPath)

	ytdlOptions := flags.Args()
	for _, name := range ytdlOptions {
		fmt.Fprintf(stdout, "Hi %s\n", name)

	}
	return nil, 0
}

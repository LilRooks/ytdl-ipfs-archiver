package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/config"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/table"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/ytdl"
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
	ytdlOptions := flags.Args()

	// Read configuration into configs variable
	err, configs := config.Parse(confPath)
	if err != nil {
		return err, errorConfig
	}

	// Configuration file is the real default
	if len(ytdlPath) == 0 {
		ytdlPath = configs.Binary
	}

	// Check binary is there
	err, ytdlPath = ytdl.ParsePath(ytdlPath)
	if err != nil {
		return err, errorYTDL
	}

	// Get the keys needed for the table
	var (
		id       string
		format   string
		location string
	)
	err, id, format = ytdl.GetIdentifiers(ytdlPath, ytdlOptions)
	if err != nil {
		return err, errorYTDL
	}

	err, location = table.Fetch(tablPath, id, format)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err, errorTable
	}

	// TODO This is basically debug stuff
	fmt.Fprintf(stdout, "Binary at %s\n", ytdlPath)
	fmt.Fprintf(stdout, "id at %s\n", id)
	fmt.Fprintf(stdout, "format at %s\n", format)
	fmt.Fprintf(stdout, "File at %s\n", location)

	for _, name := range ytdlOptions {
		fmt.Fprintf(stdout, "Hi %s\n", name)

	}
	return nil, 0
}

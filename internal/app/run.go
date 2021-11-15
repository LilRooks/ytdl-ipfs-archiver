package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/config"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/table"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/ytdl"
	"github.com/web3-storage/go-w3s-client"
)

var (
	ytdlPath string
	confPath string
	tablPath string
)

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

	_, errNotExist := os.Stat(tablPath)

	err, db := table.OpenDB(tablPath)
	if err != nil {
		return err, errorTable
	}
	defer db.Close()

	// Only initialized if the file did not originally exist
	if errors.Is(errNotExist, os.ErrNotExist) {
		fmt.Printf("[sqlite] %s doesn't exist, initializing...\n", tablPath)
		errInit := table.InitializeTable(db)
		if errInit != nil {
			return errInit, errorTable
		}
	}

	err, location = table.Fetch(db, id, format)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err, errorTable
	}

	if len(location) == 0 {
		err, filename := ytdl.Download(ytdlPath, ytdlOptions)
		if err != nil {
			return err, errorYTDL
		}
		if configs.Token == "" {
			return err, errorIPFS
		}
		c, _ := w3s.NewClient(w3s.WithToken(configs.Token))
		f, _ := os.Open(filename)

		cid, _ := c.Put(context.Background(), f)
		fmt.Printf("https://%v.ipfs.dweb.link\n", cid)
	}

	for _, name := range ytdlOptions {
		fmt.Fprintf(stdout, "Hi %s\n", name)

	}
	return nil, 0
}
